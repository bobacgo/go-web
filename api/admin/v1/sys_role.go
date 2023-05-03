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

type RoleApi struct{}

var roleService service.IRole = new(service.Role)

func (api RoleApi) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.SysRole](ctx)
	if !ok {
		return
	}
	pageResp, err := roleService.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (api RoleApi) Details(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	user, err := roleService.Details(req.ID)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, user)
}

func (api RoleApi) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RoleCreateReq](ctx)
	if !ok {
		return
	}
	if err := roleService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessCreate(ctx)
}

func (api RoleApi) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.RoleUpdateReq](ctx)
	if !ok {
		return
	}
	if err := roleService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessUpdate(ctx)
	return
}

func (api RoleApi) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := roleService.Delete(req.ID); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}