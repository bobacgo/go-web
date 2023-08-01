package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"github.com/gogoclouds/gogo/web/r"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(svc service.IUserService) *UserHandler {
	return &UserHandler{userService: svc}
}

func (h *UserHandler) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserPageQuery](ctx)
	if !ok {
		return
	}
	pageResp, err := h.userService.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (h *UserHandler) Details(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	user, err := h.userService.Details(req.ID)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, user)
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserCreateReq](ctx)
	if !ok {
		return
	}
	if err := h.userService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdateReq](ctx)
	if !ok {
		return
	}
	if err := h.userService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdateStatusReq](c)
	if !ok {
		return
	}
	if err := h.userService.UpdateStatus(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(c, err.Text)
		return
	}
	reply.SuccessUpdate(c)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	req, ok := valid.ShouldBind[model.UserUpdatePasswdReq](c)
	if !ok {
		return
	}
	if err := h.userService.UpdatePassword(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(c, err.Text)
		return
	}
	reply.SuccessUpdate(c)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := h.userService.Delete(req.ID); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}