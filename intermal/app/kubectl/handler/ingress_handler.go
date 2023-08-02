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

type IngressHandler struct {
	ingressService service.IIngressService
}

func NewIngressHandler(ingressService service.IIngressService) *IngressHandler {
	return &IngressHandler{ingressService: ingressService}
}

func (h *IngressHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Ingress](ctx)
	if !ok {
		return
	}
	if err := h.ingressService.Create(ctx.Request.Context(), req); err != nil {
		logger.Errorw("create ingress:", req.Name, err.Error())
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *IngressHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	if err := h.ingressService.Delete(ctx.Request.Context(), req); err != nil {
		logger.Errorw("delete ingress", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailDelete)
		return
	}
	reply.SuccessDelete(ctx)
}

func (h *IngressHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Ingress](ctx)
	if !ok {
		return
	}
	if err := h.ingressService.Update(ctx.Request.Context(), req); err != nil {
		logger.Errorw("update ingress", req.Name, err)
		reply.FailMsg(ctx, r.FailUpdate)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *IngressHandler) Get(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	info, err := h.ingressService.Get(ctx.Request.Context(), req)
	if err != nil {
		logger.Errorw("get ingress info", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, info)
}

func (h *IngressHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	list, err := h.ingressService.List(ctx.Request.Context(), req.Namespace)
	if err != nil {
		logger.Errorw("get ingress list", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, list)
}