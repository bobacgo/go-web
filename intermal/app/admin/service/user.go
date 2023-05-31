package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/pkg/util"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"

	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
)

type IUserService interface {
	PageList(req model.UserPageQuery) (*r.PageAnyResp, *g.Error)
	Details(id string) (model.SysUser, *g.Error)
	Create(req model.UserCreateReq) *g.Error
	Updates(q model.UserUpdateReq) *g.Error
	UpdateStatus(req model.UserUpdateStatusReq) *g.Error
	UpdatePassword(req model.UserUpdatePasswdReq) *g.Error
	Delete(ID string) *g.Error
}

var userService = new(UserService)

type UserService struct {
	g.FindByIDService[model.SysUser]
	g.UniqueService[model.SysUser]
}

func (svc *UserService) PageList(req model.UserPageQuery) (*r.PageAnyResp, *g.Error) {

	// 要求
	// 1.支持 账号名 模糊搜索
	// 2.支持 手机号 模糊搜索
	// 3.支持 昵称 模糊搜索
	// 4.支持通过角色ID查找用户

	type UserInfo struct {
		model.SysUser
		RoleInfo model.SimpleRole `json:"roleInfo" gorm:"foreignKey:RoleID;references:ID"`
	}

	db := g.DB.Model(&model.SysUser{})
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
	data, err := orm.PageAnyFind[UserInfo](db, req.PageInfo)
	return data, g.WrapError(err, r.FailRead)
}

func (svc *UserService) Details(id string) (model.SysUser, *g.Error) {
	return svc.FindByID(g.DB, id)
}

func (svc *UserService) findByUsername(tx *gorm.DB, username string) (model.SysUser, *g.Error) {
	var u model.SysUser
	err := tx.Where(&model.SysUser{Username: username}).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u, g.NewError("用户" + r.FailRecordNotFound)
	}
	return u, g.WrapError(err, "获取用户数据出错")
}

func (svc *UserService) Create(q model.UserCreateReq) *g.Error {
	var u model.SysUser
	copier.Copy(&u, &q)

	u.Username = strings.Trim(u.Username, " ")
	u.Nickname = strings.Trim(u.Nickname, " ")
	u.Password = passwordHandler.bcryptHash(q.Password)
	u.Status = enum.UserStatusEnable

	// 要求
	// 1。username、phone、email 唯一
	err := g.DB.Transaction(func(tx *gorm.DB) error {
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

func (svc *UserService) Updates(q model.UserUpdateReq) *g.Error {
	q.Username = strings.Trim(q.Username, " ")
	q.Nickname = strings.Trim(q.Nickname, " ")

	// 要求
	// 1. 传入ID查找数据是否存在
	// 2。username、phone、email 唯一
	err := g.DB.Transaction(func(tx *gorm.DB) error {
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
		return err
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailUpdate)
}

func (svc *UserService) UpdateStatus(req model.UserUpdateStatusReq) *g.Error {
	res := g.DB.Model(&model.SysUser{}).Where("id = ?", req.ID).Update("status", req.Status)
	if res.Error != nil {
		return g.WrapError(res.Error, r.FailUpdate)
	} else if res.RowsAffected == 0 {
		return g.NewError(r.FailRecordNotFound)
	}
	// TODO 退出 token
	return nil
}

func (svc *UserService) UpdatePassword(req model.UserUpdatePasswdReq) *g.Error {
	u, gerr := svc.FindByID(g.DB, req.ID)
	if gerr != nil {
		return gerr
	}

	if !passwordHandler.bcryptVerify(u.ID, req.OldPassword, u.Password) {
		return g.NewError("旧密码验证不通过")
	}

	password := passwordHandler.bcryptHash(req.NewPassword)
	err := g.DB.Model(&model.SysUser{}).Where("id = ?", req.ID).Update("password", password).Error
	// TODO 退出 token
	return g.WrapError(err, r.FailUpdate)
}

func (svc *UserService) Delete(ID string) *g.Error {
	if res := g.DB.Where("id = ?", ID).Delete(&model.SysUser{}); res.Error != nil {
		return g.WrapError(res.Error, r.FailDelete)
	} else if res.RowsAffected == 0 {
		return g.NewError(r.FailRecordNotFound)
	}
	return nil
}

func (svc *UserService) IsExist(q model.UniqueVerifyReq) (bool, *g.Error) {
	var m map[string]any
	copier.Copy(&m, &q)
	if gerr := svc.UniqueService.Verify(g.DB, m); gerr == nil {
		return false, nil
	} else if gerr.Is(g.ErrRecordRepeat) {
		return true, gerr
	} else { // 查询出错
		return true, gerr
	}
}

// ===============================================
// Password handler
// ===============================================

var passwordHandler = new(passwdHandler)

type passwdHandler struct{}

func (h *passwdHandler) bcryptVerify(userID, password, dbPassword string) bool {
	up := strings.SplitN(dbPassword, "$", 2)
	if len(up) != 2 {
		logger.Errorf("password splitN len %d", len(up))
		return false
	}
	if !util.BcryptVerify(up[0], up[1], password) {
		h.errCount(userID)
		return false
	}
	h.delErrIncr(userID)
	return true
}

func (h *passwdHandler) bcryptHash(passwd string) string {
	hash, salt := util.BcryptHash(passwd)
	return salt + hash
}

func (h *passwdHandler) errCount(userID string) {
	count, err := g.CacheDB.Incr(context.Background(), h.key(userID)).Result()
	if err != nil {
		logger.Error(userID, err)
	}
	if count > 5 {
		if err := userService.UpdateStatus(model.UserUpdateStatusReq{
			IdReq: r.IdReq{ID: userID}, Status: enum.UserStatusDisable,
		}); err != nil {
			logger.Errorf("disable user [%s] err: %v", userID, err)
		}
		h.delErrIncr(userID)
	}
}

func (h *passwdHandler) delErrIncr(userID string) {
	if err := g.CacheDB.Del(context.Background(), h.key(userID)).Err(); err != nil {
		logger.Error(userID, err)
	}
}

func (h *passwdHandler) key(userID string) string {
	return fmt.Sprintf(common.RedisKeyPasswdErrIncrFmt, userID)
}