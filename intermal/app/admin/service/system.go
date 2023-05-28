package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
)

type ISystem interface {
	Login() ([]*model.SysMenu, *g.Error)
	Logout()
	Captcha()
}

var menuService IMenuService = new(MenuService)

type System struct{}

func (s System) Login() ([]*model.SysMenu, *g.Error) {
	tree, gErr := menuService.TreeByRole("a70e443047f0403094d406b5a3c78880")
	return tree, gErr
}

func (s System) Logout() {
	//TODO implement me
	panic("implement me")
}

func (s System) Captcha() {
	//TODO implement me
	panic("implement me")
}