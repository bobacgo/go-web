package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
	"path/filepath"
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

func (api SystemApi) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		reply.FailMsg(c, "打开文件失败")
		return
	}
	filename := filepath.Base(file.Filename)
	if err = c.SaveUploadedFile(file, "./temp/"+filename); err != nil {
		logger.Error(err)
		reply.FailMsg(c, r.FailSave)
		return
	}
	reply.SuccessMsg(c, r.OKUpload)
}
