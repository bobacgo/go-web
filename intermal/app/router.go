package app

import (
	"github.com/gin-gonic/gin"
	admin_v1 "github.com/gogoclouds/go-web/api/admin/v1"
)

func loadRouter(e *gin.Engine) {
	g := e.Group("v1")

	// system
	base := g.Group("base")
	base.POST("login", admin_v1.SystemApi.Login)
	base.GET("logout", admin_v1.SystemApi.Logout)
	base.GET("captcha", admin_v1.SystemApi.Captcha)

	// sys user
	// 获取用户列表、获取用户详情、创建用户、更新用户、更新状态、更新密码、删除用户
	user := g.Group("user")
	user.POST("details", admin_v1.SysUserApi.Details)
	user.POST("", admin_v1.SysUserApi.Create)
	user.PUT("", admin_v1.SysUserApi.Update)
	user.DELETE("", admin_v1.SysUserApi.Delete)
	user.PUT("updateStatus", admin_v1.SysUserApi.UpdateStatus)
	user.PUT("updatePassword", admin_v1.SysUserApi.UpdatePassword)
	user.POST("pageList", admin_v1.SysUserApi.PageList)
	// sys role
	// sys menu
	// sys dictionary
}