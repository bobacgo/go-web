package app

import (
	"github.com/gin-gonic/gin"
	admin_v1 "github.com/gogoclouds/go-web/api/admin/v1"
	admin_v2 "github.com/gogoclouds/go-web/api/admin/v2"
)

func loadRouter(e *gin.Engine) {
	v1 := e.Group("v1")
	v2 := e.Group("v2")

	// sys menu
	menuApi_v1 := new(admin_v1.MenuApi)

	menu := v1.Group("menu")
	menu.POST("", menuApi_v1.Create)
	menu.PUT("", menuApi_v1.Update)
	menu.DELETE("", menuApi_v1.Delete)
	menu.POST("tree", menuApi_v1.Tree)
	menu.POST("simpleTree", menuApi_v1.SimpleTree)

	// sys role
	roleApi_v1 := new(admin_v1.RoleApi)

	role := v1.Group("role")
	role.POST("details", roleApi_v1.Details)
	role.POST("", roleApi_v1.Create)
	role.PUT("", roleApi_v1.Update)
	role.DELETE("", roleApi_v1.Delete)
	role.POST("pageList", roleApi_v1.PageList)
	role.POST("list", roleApi_v1.List)

	// sys user
	userApi_v1 := new(admin_v1.UserApi)

	// 获取用户列表、获取用户详情、创建用户、更新用户、更新状态、更新密码、删除用户
	user := v1.Group("user")
	user.POST("details", userApi_v1.Details)
	user.POST("", userApi_v1.Create)
	user.PUT("", userApi_v1.Update)
	user.DELETE("", userApi_v1.Delete)
	user.PUT("updateStatus", userApi_v1.UpdateStatus)
	user.PUT("updatePassword", userApi_v1.UpdatePassword)
	user.POST("pageList", userApi_v1.PageList)

	// system
	systemApi_v1 := new(admin_v1.SystemApi)

	system := v1.Group("system")
	system.POST("login", systemApi_v1.Login)
	system.GET("logout", systemApi_v1.Logout)
	system.GET("captcha", systemApi_v1.Captcha)
	system.POST("upload", systemApi_v1.Upload)

	// system v2
	systemApi_v2 := new(admin_v2.SystemApi)

	systemV2 := v2.Group("system")
	systemV2.GET("captcha", systemApi_v2.Captcha)

	// sys dictionary
}