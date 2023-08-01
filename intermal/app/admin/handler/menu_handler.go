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

type MenuHandler struct {
	menuService service.IMenuService
}

func NewMenuHandler(svc service.IMenuService) *MenuHandler {
	return &MenuHandler{
		menuService: svc,
	}
}

func (h *MenuHandler) Tree(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.MenuTreeReq](ctx)
	if !ok {
		return
	}
	pageResp, err := h.menuService.Tree(req)
	if err != nil {
		logger.Error(err)
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (h *MenuHandler) SimpleTree(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.MenuTreeReq](ctx)
	if !ok {
		return
	}
	pageResp, err := h.menuService.SimpleTree(req)
	if err != nil {
		logger.Error(err)
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (h *MenuHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.MenuCreateReq](ctx)
	if !ok {
		return
	}
	if err := h.menuService.Create(req); err != nil {
		logger.Errorf("%+v", err)
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *MenuHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.MenuUpdateReq](ctx)
	if !ok {
		return
	}
	if err := h.menuService.Save(req); err != nil {
		logger.Error(err)
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *MenuHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := h.menuService.Delete(req.ID); err != nil {
		logger.Error(err)
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}