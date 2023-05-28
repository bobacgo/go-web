package model

import (
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
)

// SysRole 系统角色
type SysRole struct {
	orm.Model
	Name        string     `json:"name" gorm:"type:varchar(100);index"`
	Description string     `json:"description" gorm:"type:varchar(1000)"`
	Menus       []*SysMenu `json:"menus" gorm:"many2many:sys_role_sys_menus"`
}

func (*SysRole) TableName() string {
	return "sys_roles"
}

type SimpleRole struct {
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"type:varchar(100)"`
}

func (*SimpleRole) TableName() string {
	return new(SysRole).TableName()
}

type RolePageListReq struct {
	r.PageInfo
	Name string `json:"name"`
}

type RoleListReq struct {
	Name string `json:"name"`
}
type RoleCreateReq struct {
	Name        string     `json:"name" binding:"required,lte=20"`
	Description string     `json:"description" binding:"lte=300"`
	Menus       []*r.IdReq `json:"menus" gorm:"many2many:sys_role_sys_menus"`
}

type RoleUpdateReq struct {
	r.IdReq
	RoleCreateReq
}