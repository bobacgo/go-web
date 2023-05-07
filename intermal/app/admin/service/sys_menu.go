package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/dao"
	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IMenu interface {
	Tree(req model.MenuTreeReq) ([]model.SysMenu, *g.Error)
	Create(req model.MenuCreateReq) error
	Updates(req model.MenuUpdateReq) *g.Error
	Delete(id string) *g.Error
}

var menuDao dao.IMenu = new(dao.Menu)

type Menu struct{}

func (svc Menu) Tree(req model.MenuTreeReq) ([]model.SysMenu, *g.Error) {
	return nil, nil
}

func (svc Menu) Create(req model.MenuCreateReq) error {
	var m model.SysMenu
	copier.Copy(&m, &req)
	return menuDao.Create(m)
}

func (svc Menu) Updates(req model.MenuUpdateReq) *g.Error {
	var m model.SysMenu
	copier.Copy(&m, &req)
	return menuDao.Updates(m)
}

func (svc Menu) Delete(ID string) *g.Error {

	// 要求：
	// 1.数据是否存在
	// 2.存在子菜单不能删除
	// 3.删除对应的角色菜单关联
	// 4.如果是按钮删除数据权限表

	err := g.DB.Transaction(func(tx *gorm.DB) error {
		var dbMenu model.SysMenu // 只返回了菜单类型值
		if err := tx.Select("menu_type").Where("id = ?", ID).Take(&dbMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return g.NewError(r.FailRecordNotFound)
			}
			return g.WrapError(err, "查找菜单信息失败")
		}
		var dbChildMenu r.IdReq
		if err := tx.Model(&model.SysMenu{}).Where(&model.SysMenu{ParentId: ID}).Take(&dbChildMenu).Error; err == nil {
			return g.NewErrorf("该菜单存在子菜单[%s]不能删除", dbChildMenu.ID)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return g.WrapError(err, "检查是否有子菜单出错")
		}

		//if err := roleOtmMenuDao.DeleteByMenuID(tx, ID); err != nil {
		//	return g.WrapError(err, "删除角色包含的菜单失败")
		//}

		if dbMenu.MenuType == enum.MenuType_Btn {
			//if err := tx.Where(&model.RoleOtmMenu{MenuID: id}).Delete(&model.RoleOtmMenu{}).Error; err != nil {
			//	return errors.WithMessage(err, "移除数据权限")
			//}
		}

		if err := tx.Where("id = ?", ID).Delete(&model.SysMenu{}).Error; err != nil {
			return g.WrapError(err, r.FailDelete)
		}
		return nil
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return nil
}