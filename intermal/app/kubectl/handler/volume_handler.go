package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/service"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"github.com/gogoclouds/gogo/web/r"
)

type VolumeHandler struct {
	volumeService service.IVolumeService
}

func NewVolumeHandler(volumeService service.IVolumeService) *VolumeHandler {
	return &VolumeHandler{
		volumeService: volumeService,
	}
}

func (h *VolumeHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Volume](ctx)
	if !ok {
		return
	}
	if err := h.volumeService.Create(ctx.Request.Context(), req); err != nil {
		logger.Errorw("create Volume:", req.Name, err.Error())
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *VolumeHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	if err := h.volumeService.Delete(ctx.Request.Context(), req); err != nil {
		logger.Errorw("delete Volume", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailDelete)
		return
	}
	reply.SuccessDelete(ctx)
}

func (h *VolumeHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Volume](ctx)
	if !ok {
		return
	}
	if err := h.volumeService.Update(ctx.Request.Context(), req); err != nil {
		logger.Errorw("update Volume", req.Name, err)
		reply.FailMsg(ctx, r.FailUpdate)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *VolumeHandler) Get(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	info, err := h.volumeService.Get(ctx.Request.Context(), req)
	if err != nil {
		logger.Errorw("get Volume info", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, info)
}

func (h *VolumeHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	list, err := h.volumeService.List(ctx.Request.Context(), req.Namespace)
	if err != nil {
		logger.Errorw("get Volume list", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, list)
}