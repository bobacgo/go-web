package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/handler"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"gorm.io/gorm"
)

// 不需要鉴权的接口
func NoAuthRouterRegister(router *gin.RouterGroup, db *gorm.DB) {
	systemHandler := handler.NewSystemHandler(service.NewSystemService(db, service.NewUserService(db), service.NewMenuService(db)))
	// 登录、获取验证码、刷新token
	router.POST("system/login", systemHandler.Login)
	router.GET("system/captcha", systemHandler.Captcha)
	router.POST("system/refreshToken", systemHandler.Refresh)
}

func RouterRegister(router *gin.RouterGroup, db *gorm.DB) {
	// sys menu
	menuService := service.NewMenuService(db)
	menuHandler := handler.NewMenuHandler(menuService)
	// 创建、删除、更新菜单、 获取详细的菜单树、获取简约的菜单树
	router.POST("admin/menu/create", menuHandler.Create)
	router.DELETE("admin/menu/delete", menuHandler.Delete)
	router.PUT("admin/menu/update", menuHandler.Update)
	router.POST("admin/menu/tree", menuHandler.Tree)
	router.POST("admin/menu/simpleTree", menuHandler.SimpleTree)

	// sys role
	roleHandler := handler.NewRoleHandler(service.NewRoleService(db))
	// 创建、删除、更新菜单、获取分页列表、不分页列表、详情角色
	router.POST("admin/role/create", roleHandler.Create)
	router.DELETE("admin/role/delete", roleHandler.Delete)
	router.PUT("admin/role/update", roleHandler.Update)
	router.POST("admin/role/pageList", roleHandler.PageList)
	router.POST("admin/role/list", roleHandler.List)
	router.POST("admin/role/details", roleHandler.Details)

	// sys user
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)
	// 获取用户列表、获取用户详情、创建用户、更新用户、更新状态、更新密码、删除用户
	router.POST("admin/user/create", userHandler.Create)
	router.DELETE("admin/user/delete", userHandler.Delete)
	router.PUT("admin/user/update", userHandler.Update)
	router.PUT("admin/user/updateStatus", userHandler.UpdateStatus)
	router.PUT("admin/user/updatePassword", userHandler.UpdatePassword)
	router.POST("admin/user/pageList", userHandler.PageList)
	router.POST("admin/user/details", userHandler.Details)

	// system
	systemHandler := handler.NewSystemHandler(service.NewSystemService(db, userService, menuService))
	// 登出、文件上传
	router.GET("admin/system/logout", systemHandler.Logout)
	router.POST("admin/system/uploadFile", systemHandler.UploadFile)

	// sys dictionary
}