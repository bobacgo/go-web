package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IMenu interface {
	Tree(req model.MenuTreeReq) ([]*model.SysMenu, *g.Error)
	Create(req model.MenuCreateReq) error
	Updates(req model.MenuUpdateReq) *g.Error
	Delete(id string) *g.Error
}

type Menu struct{}

func (svc Menu) Tree(req model.MenuTreeReq) ([]*model.SysMenu, *g.Error) {
	var menus []*model.SysMenu
	if err := g.DB.Find(&menus).Error; err != nil {
		return menus, g.WrapError(err, r.FailRead)
	}
	if len(menus) == 0 {
		return make([]*model.SysMenu, 0), nil
	}
	sameParentMap := make(map[string][]*model.SysMenu, len(menus))
	for _, m := range menus {
		cm := m
		sameParentMap[m.ParentId] = append(sameParentMap[m.ParentId], cm)
	}
	for i, m := range menus {
		children := sameParentMap[m.ID]
		if children == nil {
			children = make([]*model.SysMenu, 0)
		}
		menus[i].Children = children
	}

	// TODO 通过名称过滤树

	return sameParentMap[""], nil
}

func (svc Menu) Create(req model.MenuCreateReq) error {
	var m model.SysMenu
	copier.Copy(&m, &req)
	err := g.DB.Create(&m).Error
	return err
}

func (svc Menu) Updates(req model.MenuUpdateReq) *g.Error {
	var m model.SysMenu
	copier.Copy(&m, &req)

	// 新增菜单
	// 1.如果是按钮，更新对应的权限表
	//err := g.DB.Transaction(func(tx *gorm.DB) error {
	//	var dbMenu model.SysMenu
	//	if err := tx.Select("menu_type").Where("id = ?", o.ID).Take(&dbMenu).Error; err != nil {
	//		return errors.WithMessage(err, "获取菜单信息")
	//	}
	//	//
	//	return nil
	//})

	if res := g.DB.Updates(&m); res.Error != nil {
		return g.WrapError(res.Error, r.FailUpdate)
	} else if res.RowsAffected == 0 {
		return g.WrapError(gorm.ErrRecordNotFound, r.FailRecordNotFound)
	}
	return nil
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

		if err := tx.Where(&model.RoleOtmMenu{MenuID: ID}).Delete(&model.RoleOtmMenu{}).Error; err != nil {
			return g.WrapError(err, "删除角色菜单关联关系失败")
		}

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