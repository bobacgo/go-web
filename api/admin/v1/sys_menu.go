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

type MenuApi struct{}

var menuService service.IMenu = new(service.Menu)

func (api MenuApi) Tree(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.MenuTreeReq](ctx)
	if !ok {
		return
	}
	pageResp, err := menuService.Tree(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessRead(ctx, pageResp)
}

func (api MenuApi) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.SysMenu](ctx)
	if !ok {
		return
	}
	if err := menuService.Create(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessCreate(ctx)
}

func (api MenuApi) Update(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.SysMenu](ctx)
	if !ok {
		return
	}
	if err := menuService.Updates(req); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessUpdate(ctx)
	return
}

func (api MenuApi) Delete(ctx *gin.Context) {
	req, ok := valid.ShouldBind[r.IdReq](ctx)
	if !ok {
		return
	}
	if err := menuService.Delete(req.ID); err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessDelete(ctx)
}