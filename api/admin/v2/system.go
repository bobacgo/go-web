package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
)

type SystemApi struct{}

var systemService service.ISystem = new(service.System)

func (api *SystemApi) Captcha(c *gin.Context) {
	rsp, gErr := systemService.CaptchaV2()
	if gErr != nil {
		logger.Error(gErr.Error())
		reply.FailMsg(c, gErr.Text)
		return
	}
	reply.SuccessRead(c, rsp)
}