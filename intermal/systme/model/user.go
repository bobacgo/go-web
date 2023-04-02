package model

import (
	"github.com/gogoclouds/go-web/intermal/common"
	"time"
)

// SysUser 系统用户
type SysUser struct {
	common.OrmModel
	Username      string     `json:"username"`
	Password      string     `json:"-"`
	Nickname      string     `json:"nickname"`
	Phone         string     `json:"phone"` // 手机号
	Email         string     `json:"email"`
	Description   string     `json:"description"`
	Enable        int        `json:"enable"`        // 1.正常 2.冻结
	LastLoginTime *time.Time `json:"lastLoginTime"` // 上一次登录时间
}