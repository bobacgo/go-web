package dao

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"gorm.io/gorm"
)

type IRuleOtmMenu interface {
	DeleteByMenuID(tx *gorm.DB, ID string) error
}

type ruleOtmMenu struct{}

func (dao ruleOtmMenu) DeleteByMenuID(tx *gorm.DB, ID string) error {
	err := tx.Where(&model.RoleOtmMenu{MenuID: ID}).Delete(&model.RoleOtmMenu{}).Error
	return err
}