package model

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
)

// 设计
// 目录、菜单、按钮 使用同一个表
// path 的值是唯一

// SysMenu 系统菜单
type SysMenu struct {
	orm.Model
	ParentId string        `json:"parentId" gorm:"type:varchar(100)"`
	Name     string        `json:"name" gorm:"type:varchar(100)"`
	Path     string        `json:"path" gorm:"type:varchar(255)"`
	MenuType enum.MenuType `json:"menuType"`
	Method   string        `json:"method" gorm:"type:varchar(10)"` // net/http/method.go
	Icon     string        `json:"icon" gorm:"type:varchar(255)"`
	Sort     uint8         `json:"sort"`
	Children []*SysMenu    `json:"children" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menus"
}

type SimpleMenu struct {
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"type:varchar(100)"`
}

func (SimpleMenu) TableName() string {
	return SysMenu{}.TableName()
}

type MenuTreeReq struct {
	Name string `json:"name"`
}

type MenuCreateReq struct {
	ParentId string        `json:"parentId"`
	Name     string        `json:"name" binding:"required"`
	Path     string        `json:"path" binding:"required_if=MenuType 2"`
	MenuType enum.MenuType `json:"menuType" binding:"required,oneof=1 2 3"`
	Method   string        `json:"method" binding:"oneof='' GET POST PUT DELETE"` // net/http/method.go
	Icon     string        `json:"icon"`
	Sort     uint8         `json:"sort"`
}

type MenuUpdateReq struct {
	r.IdReq
	MenuCreateReq
}