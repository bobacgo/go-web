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

type StatefulSetHandler struct {
	statefulSetService service.IStatefulSetService
}

func NewStatefulSetHandler(statefulSetService service.IStatefulSetService) *StatefulSetHandler {
	return &StatefulSetHandler{
		statefulSetService: statefulSetService,
	}
}

func (h *StatefulSetHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.StatefulSet](ctx)
	if !ok {
		return
	}
	if err := h.statefulSetService.Create(ctx.Request.Context(), req); err != nil {
		logger.Errorw("create StatefulSet:", req.Name, err.Error())
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *StatefulSetHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	if err := h.statefulSetService.Delete(ctx.Request.Context(), req); err != nil {
		logger.Errorw("delete StatefulSet", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailDelete)
		return
	}
	reply.SuccessDelete(ctx)
}

func (h *StatefulSetHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.StatefulSet](ctx)
	if !ok {
		return
	}
	if err := h.statefulSetService.Update(ctx.Request.Context(), req); err != nil {
		logger.Errorw("update StatefulSet", req.Name, err)
		reply.FailMsg(ctx, r.FailUpdate)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *StatefulSetHandler) Get(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	info, err := h.statefulSetService.Get(ctx.Request.Context(), req)
	if err != nil {
		logger.Errorw("get StatefulSet info", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, info)
}

func (h *StatefulSetHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	list, err := h.statefulSetService.List(ctx.Request.Context(), req.Namespace)
	if err != nil {
		logger.Errorw("get StatefulSet list", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, list)
}