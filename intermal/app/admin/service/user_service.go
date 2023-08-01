package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/admin/model/enum"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/pkg/util"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"

	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
)

type IUserService interface {
	PageList(req model.UserPageQuery) (*r.PageResp[model.UserWithRole], *g.Error)
	Details(id string) (model.SysUser, *g.Error)
	Create(req model.UserCreateReq) *g.Error
	Updates(q model.UserUpdateReq) *g.Error
	UpdateStatus(req model.UserUpdateStatusReq) *g.Error
	UpdatePassword(req model.UserUpdatePasswdReq) *g.Error
	Delete(ID string) *g.Error
	FindWithRoleByUsername(*gorm.DB, string) (model.UserWithRole, *g.Error)
}

func NewUserService(db *gorm.DB) IUserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	g.FindByIDService[model.SysUser]
	g.UniqueService[model.SysUser]
	db *gorm.DB
}

func (svc *userService) PageList(req model.UserPageQuery) (*r.PageResp[model.UserWithRole], *g.Error) {

	// 要求
	// 1.支持 账号名 模糊搜索
	// 2.支持 手机号 模糊搜索
	// 3.支持 昵称 模糊搜索
	// 4.支持通过角色ID查找用户

	db := svc.db.Model(&model.SysUser{})
	if req.Username != "" {
		db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Phone != "" {
		db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Nickname != "" {
		db.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.RoleID != "" {
		db.Where("role_id = ?", req.RoleID)
	}
	db.Order("updated_at DESC").Preload("RoleInfo")
	data, err := orm.PageFind[model.UserWithRole](db, req.PageInfo)
	return data, g.WrapError(err, r.FailRead)
}

func (svc *userService) Details(id string) (model.SysUser, *g.Error) {
	return svc.FindByID(svc.db, id)
}

func (svc *userService) FindWithRoleByUsername(tx *gorm.DB, username string) (model.UserWithRole, *g.Error) {
	var u model.UserWithRole
	err := tx.Model(&model.SysUser{}).
		Where(&model.SysUser{Username: username}).
		Preload("RoleInfo").
		First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u, g.NewError("用户" + r.FailRecordNotFound)
	}
	return u, g.WrapError(err, "获取用户数据出错")
}

func (svc *userService) Create(q model.UserCreateReq) *g.Error {
	var u model.SysUser
	copier.Copy(&u, &q)

	u.Username = strings.Trim(u.Username, " ")
	u.Nickname = strings.Trim(u.Nickname, " ")
	u.Password = passwordHelper.bcryptHash(q.Password)
	u.Status = enum.UserStatusEnable

	// 要求
	// 1。username、phone、email 唯一
	err := svc.db.Transaction(func(tx *gorm.DB) error {
		m := map[string]any{
			"username": u.Username,
			"phone":    u.Phone,
			"email":    u.Email,
		}
		if gerr := svc.UniqueService.Verify(tx, m); gerr != nil {
			return gerr
		}
		err := tx.Create(&u).Error
		return err
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailCreate)
}

func (svc *userService) Updates(q model.UserUpdateReq) *g.Error {
	q.Username = strings.Trim(q.Username, " ")
	q.Nickname = strings.Trim(q.Nickname, " ")

	// 要求
	// 1. 传入ID查找数据是否存在
	// 2。username、phone、email 唯一
	err := svc.db.Transaction(func(tx *gorm.DB) error {
		u, gerr := svc.FindByID(tx, q.ID)
		if gerr != nil {
			return gerr
		}
		m := map[string]any{
			"id":       q.ID,
			"username": q.Username,
			"phone":    q.Phone,
			"email":    q.Email,
		}
		if gerr = svc.UniqueService.Verify(tx, m); gerr != nil {
			return gerr
		}
		copier.Copy(&u, &q)
		err := tx.Omit("created_at").Save(&u).Error
		if err = jwtService.Remove(u.Username); err != nil {
			return g.WrapError(err, "退出登录失败")
		}
		return err
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailUpdate)
}

func (svc *userService) UpdateStatus(req model.UserUpdateStatusReq) *g.Error {
	err := svc.db.Transaction(func(tx *gorm.DB) error {
		u, err := svc.FindByID(svc.db, req.ID)
		if err != nil {
			return g.WrapError(err, r.FailRecordNotFound)
		}
		res := svc.db.Model(&model.SysUser{}).Where("id = ?", req.ID).Update("status", req.Status)
		if res.Error != nil {
			return g.WrapError(res.Error, r.FailUpdate)
		}
		if err := jwtService.Remove(u.Username); err != nil {
			return g.WrapError(err, "退出登录失败")
		}
		return nil
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailUpdate)
}

func (svc *userService) UpdatePassword(req model.UserUpdatePasswdReq) *g.Error {
	u, gerr := svc.FindByID(svc.db, req.ID)
	if gerr != nil {
		return gerr
	}

	if !passwordHelper.bcryptVerify(u.ID, req.OldPassword, u.Password, nil) {
		return g.NewError("旧密码验证不通过")
	}

	password := passwordHelper.bcryptHash(req.NewPassword)
	err := svc.db.Model(&model.SysUser{}).Where("id = ?", req.ID).Update("password", password).Error
	if err := jwtService.Remove(u.Username); err != nil {
		logger.Errorf("修改密码，移除Token失败：%v", err)
	}
	return g.WrapError(err, r.FailUpdate)
}

func (svc *userService) Delete(ID string) *g.Error {
	if res := svc.db.Where("id = ?", ID).Delete(&model.SysUser{}); res.Error != nil {
		return g.WrapError(res.Error, r.FailDelete)
	} else if res.RowsAffected == 0 {
		return g.NewError(r.FailRecordNotFound)
	}
	return nil
}

func (svc *userService) IsExist(q model.UniqueVerifyReq) (bool, *g.Error) {
	var m map[string]any
	copier.Copy(&m, &q)
	if gerr := svc.UniqueService.Verify(svc.db, m); gerr == nil {
		return false, nil
	} else if gerr.Is(g.ErrRecordRepeat) {
		return true, gerr
	} else { // 查询出错
		return true, gerr
	}
}

// ===============================================
// Password helper
// ===============================================

var passwordHelper = new(passwdHelper)

type passwdHelper struct{}

func (h *passwdHelper) bcryptVerify(userID, hash, password string, errHandlerFunc func()) bool {
	idx := strings.LastIndex(password, "$")
	if !util.BcryptVerify(password[idx+1:], password[:idx], hash) {
		h.errCount(userID, errHandlerFunc)
		return false
	}
	h.delErrIncr(userID)
	return true
}

func (h *passwdHelper) bcryptHash(passwd string) string {
	hash, salt := util.BcryptHash(passwd)
	return hash + "$" + salt
}

func (h *passwdHelper) errCount(userID string, errHandlerFunc func()) {
	count, err := g.CacheDB.Incr(context.Background(), h.key(userID)).Result()
	if err != nil {
		logger.Error(userID, err)
	}
	if count > 5 {
		errHandlerFunc()
		h.delErrIncr(userID)
	}
}

func (h *passwdHelper) delErrIncr(userID string) {
	if err := g.CacheDB.Del(context.Background(), h.key(userID)).Err(); err != nil {
		logger.Error(userID, err)
	}
}

func (h *passwdHelper) key(userID string) string {
	return fmt.Sprintf(common.RedisKeyPasswdErrIncrFmt, userID)
}