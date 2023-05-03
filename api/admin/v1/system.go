package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
)

var SystemApi = new(systemApi)

type systemApi struct{}

var systemService service.ISystem = new(service.System)

func (s systemApi) Login(c *gin.Context) {
	systemService.Login()
}

func (s systemApi) Logout(c *gin.Context) {
	systemService.Logout()
}

func (s systemApi) Captcha(c *gin.Context) {
	systemService.Captcha()
}