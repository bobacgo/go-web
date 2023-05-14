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

type UserApi struct{}

var userService service.IUserService = new(service.UserService)

func (api UserApi) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserPageQuery](ctx)
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

func (api UserApi) Details(ctx *gin.Context) {
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

func (api UserApi) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserCreateReq](ctx)
	if !ok {
		return
	}
	if err := userService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessCreate(ctx)
}

func (api UserApi) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdateReq](ctx)
	if !ok {
		return
	}
	if err := userService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessUpdate(ctx)
	return
}

func (api UserApi) UpdateStatus(c *gin.Context) {
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

func (api UserApi) UpdatePassword(c *gin.Context) {
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

func (api UserApi) Delete(ctx *gin.Context) {
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