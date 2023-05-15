package model

import (
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
)

// SysRole 系统角色
type SysRole struct {
	orm.Model
	Name        string `json:"name" gorm:"index,type:varchar(100)"`
	Description string `json:"description" gorm:"type:varchar(1000)"`
}

func (SysRole) TableName() string {
	return "sys_roles"
}

type SimpleRole struct {
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"type:varchar(100)"`
}

func (SimpleRole) TableName() string {
	return SysRole{}.TableName()
}

// RoleOtmUser 角色用户关联(一对多)
type RoleOtmUser struct {
	RoleID string `json:"roleId"`
	UserID string `json:"userId"`
}

// RoleOtmMenu 角色菜单关联(一对多)
type RoleOtmMenu struct {
	RoleID string `json:"roleId"`
	MenuID string `json:"menuId"`
}

type RolePageListReq struct {
	r.PageInfo
	Name string `json:"name"`
}

type RoleCreateReq struct {
	Name        string   `json:"name" binding:"required,lte=20"`
	Description string   `json:"description" binding:"lte=300"`
	MenuIDs     []string `json:"menuIDs"`
	UserIDs     []string `json:"userIDs"`
}

type RoleUpdateReq struct {
	r.IdReq
	RoleCreateReq
}
