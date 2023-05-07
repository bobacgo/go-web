package model

import (
	"time"

	"github.com/gogoclouds/go-web/intermal/common"

	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
	"gorm.io/datatypes"
)

// SysUser 系统用户
type SysUser struct {
	orm.Model
	Username      string                        `json:"username" gorm:"type:varchar(100)"`
	Password      string                        `json:"-" gorm:"type:varchar(255)"`
	Nickname      string                        `json:"nickname" gorm:"type:varchar(100)"`
	Phone         string                        `json:"phone" gorm:"type:varchar(20)"` // 手机号
	RoleID        string                        `json:"roleId" gorm:"type:varchar(100)"`
	Status        enum.UserStatus               `json:"status"`
	LastLoginTime *time.Time                    `json:"lastLoginTime"` // 上一次登录时间
	Attributes    datatypes.JSONType[Attribute] `json:"attributes"`
}

type Attribute struct {
	Description string           `json:"description" gorm:"type:varchar(500)"`
	Email       string           `json:"email" gorm:"type:varchar(100)"`
	Birthday    orm.LocalTime    `json:"birthday"`
	Gender      enum.UserGenders `json:"gender"`
	Address     common.Location  `json:"address"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

type SimpleUser struct {
	ID       string `json:"id" gorm:"primarykey"`
	Nickname string `json:"nickname" gorm:"type:varchar(100)"`
}

func (SimpleUser) TableName() string {
	return SysUser{}.TableName()
}

type UserPageQuery struct {
	r.PageInfo
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type UserCreateReq struct {
	Username   string                        `json:"username" binding:"required"`
	Password   string                        `json:"password" binding:"required"` // 密码强度要求
	Nickname   string                        `json:"nickname" binding:"required"`
	Phone      string                        `json:"phone" binding:"required,number,startswith=1,len=11"` // 手机号
	Attributes datatypes.JSONType[Attribute] `json:"attributes"`
}

type UserUpdateStatusReq struct {
	r.IdReq
	Status enum.UserStatus `json:"status" binding:"required"`
}

type UserUpdatePasswdReq struct {
	r.IdReq
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
