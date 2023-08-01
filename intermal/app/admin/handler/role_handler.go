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

type RoleHandler struct {
	roleService service.IRoleService
}

func NewRoleHandler(roleService service.IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RolePageListReq](ctx)
	if !ok {
		return
	}
	pageResp, err := h.roleService.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (h *RoleHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RoleListReq](ctx)
	if !ok {
		return
	}
	list, err := h.roleService.List(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, list)
}

func (h *RoleHandler) Details(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	user, err := h.roleService.Details(req.ID)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, user)
}

func (h *RoleHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RoleCreateReq](ctx)
	if !ok {
		return
	}
	if err := h.roleService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *RoleHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RoleUpdateReq](ctx)
	if !ok {
		return
	}
	if err := h.roleService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsgDetails(ctx, err.Text, err.Misc)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *RoleHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := h.roleService.Delete(req.ID); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}