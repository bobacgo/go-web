package app

import (
	"github.com/gin-gonic/gin"
	admin_v1 "github.com/gogoclouds/go-web/api/admin/v1"
)

func loadRouter(e *gin.Engine) {
	g := e.Group("v1")

	// sys menu
	menuApi_v1 := new(admin_v1.MenuApi)

	menu := g.Group("menu")
	menu.POST("", menuApi_v1.Create)
	menu.PUT("", menuApi_v1.Update)
	menu.DELETE("", menuApi_v1.Delete)
	menu.POST("tree", menuApi_v1.Tree)
	menu.POST("simpleTree", menuApi_v1.SimpleTree)

	// sys role
	roleApi_v1 := new(admin_v1.RoleApi)

	role := g.Group("role")
	role.POST("details", roleApi_v1.Details)
	role.POST("", roleApi_v1.Create)
	role.PUT("", roleApi_v1.Update)
	role.DELETE("", roleApi_v1.Delete)
	role.POST("pageList", roleApi_v1.PageList)

	// sys user
	userApi_v1 := new(admin_v1.UserApi)

	// 获取用户列表、获取用户详情、创建用户、更新用户、更新状态、更新密码、删除用户
	user := g.Group("user")
	user.POST("details", userApi_v1.Details)
	user.POST("", userApi_v1.Create)
	user.PUT("", userApi_v1.Update)
	user.DELETE("", userApi_v1.Delete)
	user.PUT("updateStatus", userApi_v1.UpdateStatus)
	user.PUT("updatePassword", userApi_v1.UpdatePassword)
	user.POST("pageList", userApi_v1.PageList)

	// system
	systemApi_v1 := new(admin_v1.SystemApi)

	base := g.Group("base")
	base.POST("login", systemApi_v1.Login)
	base.GET("logout", systemApi_v1.Logout)
	base.GET("captcha", systemApi_v1.Captcha)
	base.POST("upload", systemApi_v1.Upload)

	// sys dictionary
}