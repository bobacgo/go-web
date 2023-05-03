package model

import "github.com/gogoclouds/gogo/web/orm"

// SysRole 系统角色
type SysRole struct {
	orm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RoleOtmUser 角色用户关联(一对多)
type RoleOtmUser struct {
	RoleID string `json:"roleId"`
	UserID string `json:"userId"`
}

// RoleMtmMenu 角色菜单关联(多对多)
type RoleMtmMenu struct {
	RoleID string `json:"roleId"`
	MenuID string `json:"menuId"`
}