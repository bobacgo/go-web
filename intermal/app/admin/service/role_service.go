package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IRoleService interface {
	PageList(q model.RolePageListReq) (*r.PageResp[model.SysRole], *g.Error)
	List(q model.RoleListReq) ([]model.SimpleRole, *g.Error)
	Details(id string) (model.SysRole, *g.Error)
	Create(q model.RoleCreateReq) *g.Error
	Updates(q model.RoleUpdateReq) *g.Error
	Delete(id string) *g.Error
}

type roleService struct {
	g.FindByIDService[model.SysRole]
	g.UniqueService[model.SysRole]
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) IRoleService {
	return &roleService{
		db: db,
	}
}

func (svc *roleService) PageList(q model.RolePageListReq) (*r.PageResp[model.SysRole], *g.Error) {

	// 要求
	// 1.name 模糊搜索
	db := svc.db.Model(&model.SysRole{})
	if q.Name != "" {
		db.Where("name LIKE ?", q.Name)
	}
	db.Order("updated_at DESC")
	pageData, err := orm.PageFind[model.SysRole](db, q.PageInfo)
	return pageData, g.WrapError(err, r.FailRead)
}

func (svc *roleService) List(q model.RoleListReq) ([]model.SimpleRole, *g.Error) {

	// 要求
	// 1.name 模糊搜索
	db := svc.db.Model(&model.SysRole{})
	if q.Name != "" {
		db.Where("name LIKE ?", q.Name)
	}
	db.Order("updated_at DESC")
	pageData, err := orm.PageFind[model.SimpleRole](db, r.PageInfo{Page: 0, PageSize: 100})
	return pageData.List, g.WrapError(err, r.FailRead)
}

func (svc *roleService) Details(id string) (model.SysRole, *g.Error) {
	var role model.SysRole
	if err := svc.db.Model(model.SysRole{}).Where("id = ?", id).
		Preload("Menus").
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role, g.NewError(r.FailRecordNotFound)
		}
		return role, g.WrapError(err, r.FailRead)
	}
	return role, nil
}

func (svc *roleService) Create(q model.RoleCreateReq) *g.Error {

	// 要求
	// 1.name 不能重复
	// 2.角色菜单关联关系数据写入

	err := svc.db.Transaction(func(tx *gorm.DB) error {
		m := map[string]any{"name": q.Name}
		if gErr := svc.UniqueService.Verify(tx, m); gErr != nil {
			return gErr
		}
		var role model.SysRole
		copier.Copy(&role, &q)
		if err := tx.Omit("Menus.*").Create(&role).Error; err != nil {
			return g.WrapError(err, "创建角色出错")
		}
		return nil
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailCreate)
}

func (svc *roleService) Updates(q model.RoleUpdateReq) *g.Error {
	// 要求
	// 1.name 不能重复
	// 2.角色菜单关联关系数据写入

	err := svc.db.Transaction(func(tx *gorm.DB) error {
		m := map[string]any{"id": q.ID, "name": q.Name}
		if _, gErr := svc.FindByID(tx, q.ID); gErr != nil {
			return gErr
		}
		if gErr := svc.UniqueService.Verify(tx, m); gErr != nil {
			return gErr
		}

		// many2many 不会清空 关联关系表，需要手动清除
		// DELETE FROM `sys_role_sys_menus` WHERE `sys_role_sys_menus`.`sys_role_id` = 'a70e443047f0403094d406b5a3c78880'
		if err := tx.Model(&model.SysRole{Model: orm.Model{ID: q.ID}}).Association("Menus").Clear(); err != nil {
			return g.WrapError(err, "删除关联出错")
		}

		var role model.SysRole
		copier.Copy(&role, &q)
		if err := tx.Omit("Menus.*").Updates(&role).Error; err != nil {
			return g.WrapError(err, "更新角色出错")
		}
		return nil
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailCreate)
}

func (svc *roleService) Delete(ID string) *g.Error {

	// 要求
	// 1.角色下有在使用用户不能删除
	// 2.删除角色下所有菜单关系

	err := svc.db.Debug().Transaction(func(tx *gorm.DB) error {
		if _, gErr := svc.FindByID(tx, ID); gErr != nil {
			return gErr
		}
		if gErr := new(userService).Verify(tx, map[string]any{"role_id": ID}); gErr != nil {
			gErr.Text = "该角色下存在用户不能删除"
			return gErr
		}

		// DELETE FROM `sys_role_sys_menus` WHERE `sys_role_sys_menus`.`sys_role_id` = 'a70e443047f0403094d406b5a3c78880'
		if err := tx.Model(&model.SysRole{Model: orm.Model{ID: ID}}).Association("Menus").Clear(); err != nil {
			return g.WrapError(err, "删除关联出错")
		}
		if err := tx.Where("id = ?", ID).Delete(&model.SysRole{}).Error; err != nil {
			return g.WrapError(err, r.FailDelete)
		}
		return nil
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailDelete)
}