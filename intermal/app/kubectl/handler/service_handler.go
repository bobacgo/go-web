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

type ServiceHandler struct {
	svcService service.ISvcService
}

func NewServiceHandler(svc service.ISvcService) *ServiceHandler {
	return &ServiceHandler{
		svcService: svc,
	}
}

func (h *ServiceHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Service](ctx)
	if !ok {
		return
	}
	if err := h.svcService.Create(ctx.Request.Context(), req); err != nil {
		logger.Errorw("create service:", req.Name, err.Error())
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *ServiceHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	if err := h.svcService.Delete(ctx.Request.Context(), req); err != nil {
		logger.Errorw("delete service", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailDelete)
		return
	}
	reply.SuccessDelete(ctx)
}

func (h *ServiceHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Service](ctx)
	if !ok {
		return
	}
	if err := h.svcService.Update(ctx.Request.Context(), req); err != nil {
		logger.Errorw("update service", req.Name, err)
		reply.FailMsg(ctx, r.FailUpdate)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *ServiceHandler) Get(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	serviceInfo, err := h.svcService.Get(ctx.Request.Context(), req)
	if err != nil {
		logger.Errorw("get service info", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, serviceInfo)
}

func (h *ServiceHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	serviceList, err := h.svcService.List(ctx.Request.Context(), req.Namespace)
	if err != nil {
		logger.Errorw("get service list", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, serviceList)
}