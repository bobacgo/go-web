package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
)

type IRole interface {
	PageList(req model.SysRole) (interface{}, *g.Error)
	Details(id string) (interface{}, *g.Error)
	Create(req model.RoleCreateReq) *g.Error
	Updates(req model.RoleUpdateReq) *g.Error
	Delete(id string) *g.Error
}

type Role struct {
}

func (svc Role) PageList(req model.SysRole) (interface{}, *g.Error) {
	//TODO implement me
	panic("implement me")
}

func (svc Role) Details(id string) (interface{}, *g.Error) {
	//TODO implement me
	panic("implement me")
}

func (svc Role) Create(req model.RoleCreateReq) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc Role) Updates(req model.RoleUpdateReq) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc Role) Delete(id string) *g.Error {
	//TODO implement me
	panic("implement me")
}
