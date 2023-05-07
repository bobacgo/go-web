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
	FindByID(tx *gorm.DB, ID string) (model.SysMenu, error)
	FindOneChildMenuByID(tx *gorm.DB, menuID string) (childID string, err error)
	Create(o model.SysMenu) error
	Updates(o model.SysMenu) *g.Error
	Delete(tx *gorm.DB, id string) error
}

type Menu struct{}

func (dao Menu) SimpleTree(q model.MenuTreeReq) ([]model.SimpleMenu, *g.Error) {
	//TODO implement me
	panic("implement me")
}

func (dao Menu) Tree(q model.MenuTreeReq) ([]model.SysMenu, *g.Error) {
	return nil, nil
}

func (dao Menu) FindByID(tx *gorm.DB, ID string) (model.SysMenu, error) {
	var m model.SysMenu
	err := tx.Where("id = ?", ID).First(&m).Error
	return m, err
}

func (dao Menu) FindOneChildMenuByID(tx *gorm.DB, menuID string) (childID string, err error) {
	var dbChildMenu r.IdReq
	if err = tx.Model(&model.SysMenu{}).Where(&model.SysMenu{ParentId: menuID}).First(&dbChildMenu).Error; err == nil {
		return dbChildMenu.ID, g.ErrDateBusy
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	return "", nil
}

func (dao Menu) Create(o model.SysMenu) error {
	err := g.DB.Create(&o).Error
	return errors.WithStack(err)
}

func (dao Menu) Updates(o model.SysMenu) *g.Error {

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
	if res := g.DB.Updates(&o); res.Error != nil {
		return g.WrapError(res.Error, r.FailUpdate)
	} else if res.RowsAffected == 0 {
		return g.WrapError(gorm.ErrRecordNotFound, r.FailRecordNotFound)
	}
	return nil
}

func (dao Menu) Delete(tx *gorm.DB, ID string) error {
	//if err := roleOtmMenuDao.DeleteByMenuID(tx, ID); err != nil {
	//	return g.WrapError(err, "删除角色菜单关联关系")
	//}
	err := tx.Where("id = ?", ID).Delete(&model.SysMenu{}).Error
	return err
}