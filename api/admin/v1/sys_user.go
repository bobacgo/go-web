package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"github.com/gogoclouds/gogo/web/r"
)

var SysUserApi = new(sysUserApi)

type sysUserApi struct{}

var userService service.IUser = new(service.User)

func (api sysUserApi) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.PageQuery](ctx)
	if !ok {
		return
	}
	pageResp, err := userService.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (api sysUserApi) Details(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	user, err := userService.Details(req.ID)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, user)
}

func (api sysUserApi) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserCreateReq](ctx)
	if !ok {
		return
	}
	if err := userService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessCreate(ctx)
}

func (api sysUserApi) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.SysUser](ctx)
	if !ok {
		return
	}
	if err := userService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessUpdate(ctx)
	return
}

func (api sysUserApi) UpdateStatus(c *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdateStatusReq](c)
	if !ok {
		return
	}
	if err := userService.UpdateStatus(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(c, err.Text)
		return
	}
	reply.SuccessUpdate(c)
}

func (api sysUserApi) UpdatePassword(c *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdatePasswdReq](c)
	if !ok {
		return
	}
	if err := userService.UpdatePassword(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(c, err.Text)
		return
	}
	reply.SuccessUpdate(c)
}

func (api sysUserApi) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := userService.Delete(req.ID); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}