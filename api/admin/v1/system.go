package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
)

type SystemApi struct{}

var systemService service.ISystem = new(service.System)

func (api SystemApi) Login(c *gin.Context) {
	systemService.Login()
}

func (api SystemApi) Logout(c *gin.Context) {
	systemService.Logout()
}

func (api SystemApi) Captcha(c *gin.Context) {
	systemService.Captcha()
}
