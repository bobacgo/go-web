package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/go-web/intermal/util"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"github.com/gogoclouds/gogo/web/r"
	"path/filepath"
)

type SystemApi struct{}

var systemService service.ISystem = new(service.SystemService)

func (api *SystemApi) Login(c *gin.Context) {
	req, ok := valid.ShouldBind[model.LoginReq](c)
	if !ok {
		return
	}
	resp, gErr := systemService.Login(req)
	if gErr != nil {
		logger.Error(gErr.Error())
		if gErr.Is(service.ErrUserDisable) {
			reply.Fail(c, 4100, gErr.Text) // TODO
		} else {
			reply.FailMsg(c, gErr.Text)
		}
		return
	}
	reply.SuccessMsgData(c, "登录成功", resp)
}

func (api *SystemApi) Logout(c *gin.Context) {
	err := systemService.Logout(util.ContextUsername(c))
	if err != nil {
		logger.Error(err)
		reply.FailMsg(c, "退出登录失败")
		return
	}
	reply.SuccessMsg(c, "退出登录失败")
}

func (api *SystemApi) Captcha(c *gin.Context) {
	captchaRsp, err := systemService.Captcha()
	if err != nil {
		logger.Error(err)
		reply.FailMsg(c, "获取验证码失败")
		return
	}
	reply.SuccessRead(c, captchaRsp)
}

func (api *SystemApi) Refresh(c *gin.Context) {
	req, ok := valid.ShouldBind[model.RefreshTokenVo](c)
	if !ok {
		return
	}
	res, err := systemService.Refresh(req)
	if err != nil {
		logger.Error(err)
		reply.FailMsg(c, "更新令牌失败")
		return
	}
	reply.SuccessData(c, res)
}

func (api *SystemApi) Upload(c *gin.Context) {
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
