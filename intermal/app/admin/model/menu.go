package model

import "github.com/gogoclouds/gogo/web/orm"

// SysMenu 系统菜单
type SysMenu struct {
	orm.Model
	ParentId string     `json:"parentId"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	MenuType int        `json:"menuType"`
	Method   string     `json:"method"` // 请求方法 GET | POST | PUT | DELETE
	Icon     string     `json:"icon"`
	Sort     int        `json:"sort"`
	Children []*SysMenu `json:"children" gorm:"-"`
}