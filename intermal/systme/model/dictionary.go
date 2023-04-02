package model

import (
	"github.com/gogoclouds/go-web/intermal/common"
)

// SysDictionary 系统字典
type SysDictionary struct {
	common.OrmModel
	Namespace   string `json:"namespace" orm:"default:default,comment:业务领域"`
	ParentCode  *uint  `json:"parentCode" orm:"primaryKey,comment:字典类型"`
	Code        *uint  `json:"code" orm:"primaryKey,comment:字典code"`
	Value       string `json:"value" orm:"not null,comment:字典value"`
	Description string `json:"description" orm:"comment:说明"`
}