package dao

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IMenu interface {
	Tree(q model.MenuTreeReq) ([]model.SysMenu, *g.Error)
	SimpleTree(q model.MenuTreeReq) ([]model.SimpleMenu, *g.Error)
	Create(o model.SysMenu) *g.Error
	Updates(o model.SysMenu) *g.Error
	Delete(id string) *g.Error
}

type Menu struct{}

func (dao Menu) SimpleTree(q model.MenuTreeReq) ([]model.SimpleMenu, *g.Error) {
	//TODO implement me
	panic("implement me")
}

func (dao Menu) Updates(o model.SysMenu) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (dao Menu) Delete(id string) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (dao Menu) Tree(q model.MenuTreeReq) ([]model.SysMenu, *g.Error) {
	return nil, nil
}

func (dao Menu) Create(o model.SysMenu) *g.Error {

	// 新增菜单
	// 目录、菜单、按钮 使用同一个表
	// path 的值是唯一

	err := g.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where(&model.SysMenu{Path: o.Path}).First(&model.SysMenu{}).Error
		if err == nil {
			return g.ErrRecordRepeat
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Create(&o).Error
		}
		return err
	})
	if errors.Is(err, g.ErrRecordRepeat) {
		return g.NewErrorf("[%s]: Path已经存在", o.Path)
	}
	return g.WrapError(err, r.FailCreate)
}