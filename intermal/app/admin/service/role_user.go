package service

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"gorm.io/gorm"
)

var RoleUserService = new(roleUserService)

type roleUserService struct{}

func (svc *roleUserService) findByRoleID(tx *gorm.DB, roleID string) ([]model.RoleOtmUser, *g.Error) {
	var rms []model.RoleOtmUser
	err := tx.Where("role_id = ?", roleID).Find(&rms).Error
	return rms, g.WrapError(err, "获取角色下的用户失败")
}

func (svc *roleUserService) createInBatches(tx *gorm.DB, roleID string, userIDs []string) error {
	var rus []model.RoleOtmUser
	for _, userID := range userIDs {
		rus = append(rus, model.RoleOtmUser{RoleID: roleID, UserID: userID})
	}
	err := tx.CreateInBatches(rus, 1000).Error
	return err
}
