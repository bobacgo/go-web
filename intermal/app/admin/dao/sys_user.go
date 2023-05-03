package dao

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

type IUser interface {
	PageList(q model.UserPageQuery) (*r.PageAnyResp, *g.Error)
	FindByID(ID string) (model.SysUser, *g.Error)
	Create(u model.SysUser) *g.Error
	Update(u model.SysUser) *g.Error
	Delete(id string) *g.Error
}

type User struct{}

func (dao *User) PageList(q model.UserPageQuery) (*r.PageAnyResp, *g.Error) {
	// 要求
	// 1.支持 账号名 模糊搜索
	// 2.支持 手机号 模糊搜索
	// 3.支持 昵称 模糊搜索

	type UserInfo struct {
		model.SysUser
		RoleInfo model.SimpleRole `json:"roleInfo" gorm:"foreignKey:RoleID;references:ID"`
	}

	db := g.DB.Model(&model.SysUser{})
	if q.Username != "" {
		db.Where("username LIKE ?", "%"+q.Username+"%")
	}
	if q.Phone != "" {
		db.Where("phone LIKE ?", "%"+q.Phone+"%")
	}
	if q.Nickname != "" {
		db.Where("nickname LIKE ?", "%"+q.Nickname+"%")
	}
	db.Order("updated_at DESC").Preload("RoleInfo")
	data, err := orm.PageAnyFind[UserInfo](db, q.PageInfo)
	return data, g.WrapError(err, r.FailRead)
}

func (dao *User) FindByID(id string) (res model.SysUser, gerr *g.Error) {
	err := g.DB.Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		gerr = g.NewError(r.FailRecordNotFound)
		return
	}
	gerr = g.WrapError(err, r.FailRead)
	return
}

func (dao *User) Create(u model.SysUser) *g.Error {

	// 要求
	// 1。username 唯一
	// 2. phone 唯一

	err := g.DB.Transaction(func(tx *gorm.DB) error {
		var dbList []model.SysUser
		if err := tx.Where("username = ? or phone = ?", u.Username, u.Phone).Find(&dbList).Error; err != nil {
			return err
		}
		if len(dbList) != 0 {
			errMsg := ""
			for _, v := range dbList {
				if v.Username == u.Username {
					errMsg += "|" + u.Username
				}
				if v.Phone == u.Phone {
					errMsg += "|" + u.Phone
				}
			}
			return errors.WithMessage(g.ErrRecordRepeat, strings.TrimLeft(errMsg, "|"))
		}
		err := tx.Create(&u).Error
		return err
	})
	if errors.Is(err, g.ErrRecordRepeat) {
		return g.NewError(err.Error())
	}
	return g.WrapError(err, r.FailCreate)
}

func (dao *User) Update(u model.SysUser) *g.Error {

	// 要求
	// 1。username 唯一
	// 2. phone 唯一

	if u.ID == "" {
		return g.NewError(r.FailIDNotNil)
	}
	err := g.DB.Transaction(func(tx *gorm.DB) error {
		var dbList []model.SysUser
		if err := tx.Where("id = ? or username = ? or phone = ?", u.ID, u.Username, u.Phone).Find(&dbList).Error; err != nil {
			return err
		}
		if len(dbList) >= 2 {
			errMsg := ""
			for _, v := range dbList {
				if v.ID != u.ID && v.Username == u.Username {
					errMsg += "|" + u.Username
				}
				if v.ID != u.ID && v.Phone == u.Phone {
					errMsg += "|" + u.Phone
				}
			}
			return errors.WithMessage(g.ErrRecordRepeat, strings.TrimLeft(errMsg, "|"))
		}
		if len(dbList) == 0 || (len(dbList) == 1 && dbList[0].ID != u.ID) {
			return gorm.ErrRecordNotFound
		}
		err := tx.Updates(&u).Error
		return err
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return g.NewError(r.FailRecordNotFound)
	}
	if errors.Is(err, g.ErrRecordRepeat) {
		return g.NewErrorf(err.Error())
	}
	return g.WrapError(err, r.FailUpdate)
}

func (dao *User) Delete(ID string) *g.Error {
	res := g.DB.Delete(&model.SysUser{}, ID)
	if res.RowsAffected == 0 {
		return g.NewError(r.FailRecordNotFound)
	}
	return g.WrapError(res.Error, r.FailDelete)
}