package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/dao"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
)

type IMenu interface {
	Tree(req model.MenuTreeReq) ([]model.SysMenu, *g.Error)
	Create(req model.SysMenu) *g.Error
	Updates(req model.SysMenu) *g.Error
	Delete(id string) *g.Error
}

var menuDao dao.IMenu = new(dao.Menu)

type Menu struct{}

func (svc Menu) Updates(req model.SysMenu) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc Menu) Delete(id string) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc Menu) Tree(req model.MenuTreeReq) ([]model.SysMenu, *g.Error) {
	return nil, nil
}

func (svc Menu) Create(req model.SysMenu) *g.Error {
	return menuDao.Create(req)
}