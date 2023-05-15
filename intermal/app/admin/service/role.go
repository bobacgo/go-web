package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/orm"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type IRole interface {
	PageList(q model.RolePageListReq) (*r.PageResp[model.SysRole], *g.Error)
	Details(id string) (model.RoleUpdateReq, *g.Error)
	Create(q model.RoleCreateReq) *g.Error
	Updates(q model.RoleUpdateReq) *g.Error
	Delete(id string) *g.Error
}

type Role struct {
	g.FindByIDService[model.SysRole]
	g.UniqueService[model.SysRole]
}

func (svc Role) PageList(q model.RolePageListReq) (*r.PageResp[model.SysRole], *g.Error) {

	// 要求
	// 1.name 模糊搜索
	db := g.DB.Model(&model.SysRole{})
	if q.Name != "" {
		db.Where("name LIKE ?", q.Name)
	}
	db.Order("updated_at DESC")
	pageData, err := orm.PageFind[model.SysRole](db, q.PageInfo)
	return pageData, g.WrapError(err, r.FailRead)
}

func (svc Role) Details(id string) (model.RoleUpdateReq, *g.Error) {
	// 要求
	// 1.name 不能重复
	// 2.角色菜单关联关系数据清空覆盖
	// 3.角色用户关联关系数据清空覆盖
	var roleInfo model.RoleUpdateReq
	return roleInfo, nil
}

func (svc Role) Create(q model.RoleCreateReq) *g.Error {

	// 要求
	// 1.name 不能重复
	// 2.角色菜单关联关系数据写入
	// 3.角色用户关联关系表

	err := g.DB.Transaction(func(tx *gorm.DB) error {
		m := map[string]any{"name": q.Name}
		if gErr := svc.UniqueService.Verify(tx, m); gErr != nil {
			return gErr
		}
		var role model.SysRole
		copier.Copy(&role, &q)
		if err := tx.Create(&role).Error; err != nil {
			return g.WrapError(err, "创建角色出错")
		}
		if len(q.MenuIDs) > 0 {
			if err := RoleMenuService.createInBatches(tx, role.ID, q.MenuIDs); err != nil {
				logger.Error("角色菜单关联关系数据写入失败", err)
			}
		}
		if len(q.UserIDs) > 0 {
			if err := RoleUserService.createInBatches(tx, role.ID, q.UserIDs); err != nil {
				logger.Error("角色用户关联关系数据写入失败", err)
			}
		}
		return nil

	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailCreate)
}

func (svc Role) Updates(q model.RoleUpdateReq) *g.Error {
	// 要求
	// 1.name 不能重复
	// 2.角色菜单关联关系数据清空覆盖
	// 3.角色用户关联关系数据清空覆盖
	return nil
}

func (svc Role) Delete(ID string) *g.Error {

	// 要求
	// 1.角色下有在使用用户不能删除
	// 2.删除角色下所有菜单关系

	err := g.DB.Transaction(func(tx *gorm.DB) error {
		if _, gErr := svc.FindByID(tx, ID); gErr != nil {
			return gErr
		}

		rus, gErr := RoleUserService.findByRoleID(tx, ID)
		if gErr != nil {
			return gErr
		}
		if len(rus) > 0 {
			return g.WrapError(g.ErrDateBusy, "该角色下存在用户不能删除")
		}

		if err := tx.Where("id = ?", ID).Delete(&model.SysRole{}).Error; err != nil {
			return g.WrapError(err, r.FailDelete)
		}

		gErr = RoleMenuService.deleteByRoleID(tx, ID)
		return gErr
	})
	if gerr, ok := err.(*g.Error); ok {
		return gerr
	}
	return g.WrapError(err, r.FailDelete)
}
