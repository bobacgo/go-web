package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"gorm.io/gorm"
)

var RoleMenuService = new(roleMenuService)

type roleMenuService struct{}

func (svc *roleMenuService) findByRoleID(tx *gorm.DB, roleID string) ([]model.RoleOtmMenu, *g.Error) {
	var rms []model.RoleOtmMenu
	err := tx.Where("role_id = ?", roleID).Find(&rms).Error
	return rms, g.WrapError(err, "获取角色下的菜单失败")
}

func (svc *roleMenuService) deleteByMenuID(tx *gorm.DB, menuID string) *g.Error {
	err := tx.Where("menu_id = ?", menuID).Delete(&model.RoleOtmMenu{}).Error
	return g.WrapError(err, "删除指定角色下的菜单失败")
}