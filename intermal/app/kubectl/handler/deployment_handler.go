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

type DeploymentHandler struct {
	deploymentService service.IDeploymentService
}

func NewDeploymentHandler(svc service.IDeploymentService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: svc,
	}
}

func (h *DeploymentHandler) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Deployment](ctx)
	if !ok {
		return
	}
	if err := h.deploymentService.Create(ctx.Request.Context(), req); err != nil {
		logger.Errorw("create deployment:", req.Name, err.Error())
		reply.FailMsg(ctx, r.FailCreate)
		return
	}
	reply.SuccessCreate(ctx)
}

func (h *DeploymentHandler) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	if err := h.deploymentService.Delete(ctx.Request.Context(), req); err != nil {
		logger.Errorw("delete deployment", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailDelete)
		return
	}
	reply.SuccessDelete(ctx)
}

func (h *DeploymentHandler) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.Deployment](ctx)
	if !ok {
		return
	}
	if err := h.deploymentService.Update(ctx.Request.Context(), req); err != nil {
		logger.Errorw("update deployment", req.Name, err)
		reply.FailMsg(ctx, r.FailUpdate)
		return
	}
	reply.SuccessUpdate(ctx)
}

func (h *DeploymentHandler) List(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	deploymentList, err := h.deploymentService.List(ctx.Request.Context(), req.Namespace)
	if err != nil {
		logger.Errorw("get deployment list", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, deploymentList)
}

func (h *DeploymentHandler) Get(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.NamespaceWithName](ctx)
	if !ok {
		return
	}
	deploymentInfo, err := h.deploymentService.Get(ctx.Request.Context(), req)
	if err != nil {
		logger.Errorw("get deployment info", req.Namespace, req.Name, "errMsg", err.Error())
		reply.FailMsg(ctx, r.FailRead)
		return
	}
	reply.SuccessData(ctx, deploymentInfo)
}