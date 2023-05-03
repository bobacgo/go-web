package model

import (
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
)

// SysRole 系统角色
type SysRole struct {
	orm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
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

type RoleCreateReq struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MenuIDs     []string `json:"menuIDs"`
}

type RoleUpdateReq struct {
	r.IdReq
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MenuIDs     []string `json:"menuIDs"`
}