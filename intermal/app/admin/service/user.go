package service

import (
	"github.com/gogoclouds/gogo/web/orm"
	"gorm.io/gorm"
	"strings"

	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/pkg/util"
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

func (svc *UserService) Create(q model.UserCreateReq) *g.Error {
	var u model.SysUser
	copier.Copy(&u, &q)

	u.Username = strings.Trim(u.Username, " ")
	u.Nickname = strings.Trim(u.Nickname, " ")
	u.Password = svc.bcrypt().Hash(q.Password)
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
	return nil
}

func (svc *UserService) UpdatePassword(req model.UserUpdatePasswdReq) *g.Error {
	u, gerr := svc.FindByID(g.DB, req.ID)
	if gerr != nil {
		return gerr
	}

	bcrypt := svc.bcrypt()
	if !bcrypt.Check(u.Password, req.OldPassword) {
		return g.NewError("旧密码验证不通过")
	}

	password := bcrypt.Hash(req.NewPassword)
	err := g.DB.Model(&model.SysUser{}).Where("id = ?", req.ID).Update("password", password).Error
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

func (svc *UserService) bcrypt() util.Bcrypt {
	serviceKV := g.Conf.AppServiceKV()
	bcrypt := util.NewBcrypt(serviceKV["salt"].(string))
	return bcrypt
}