package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

type ICasbin interface {
	Casbin() *casbin.SyncedCachedEnforcer
}

type Casbin struct{}

var (
	once                 sync.Once
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
)

func (svc *Casbin) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		var err error
		if syncedCachedEnforcer, err = svc.casbin(); err != nil {
			logger.Error(err)
			return
		}
	})
	return syncedCachedEnforcer
}

func (svc *Casbin) RemoveFilteredPolicy(v int, p ...string) bool {
	e := svc.Casbin()
	ok, err := e.RemoveFilteredPolicy(v, p...)
	if err != nil {
		logger.Error(err)
	}
	return ok
}

func (svc *Casbin) UpdateCasbinApi(menuID string, path, method string) error {
	//casbinDao.UpdateByMenuID(menuID, path, method)
	e := svc.Casbin()
	err := e.InvalidateCache()
	return err
}

func (svc *Casbin) casbin() (*casbin.SyncedCachedEnforcer, error) {
	a, err := gormadapter.NewAdapterByDB(g.DB)
	if err != nil {
		return nil, errors.WithMessage(err, "适配数据库失败请检查casbin表是否为InnoDB引擎!")
	}
	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, errors.WithMessage(err, "Casbin Model 加载失败")
	}
	syncedCachedEnforcer, err := casbin.NewSyncedCachedEnforcer(m, a)
	syncedCachedEnforcer.SetExpireTime(60 * 60)
	err = syncedCachedEnforcer.LoadPolicy()
	return syncedCachedEnforcer, err
}

func (svc *Casbin) UpdateByMenuID(tx *gorm.DB, menuID, path, method string) *g.Error {
	cr := gormadapter.CasbinRule{
		V1: path,
		V2: method,
	}
	if res := tx.Model(&gormadapter.CasbinRule{}).Where("v3 = ?", menuID).Updates(&cr); res.Error != nil {
		return g.WrapError(res.Error, "更新Casbin失败")
	} else if res.RowsAffected == 0 {
		return g.WrapError(gorm.ErrRecordNotFound, r.FailRecordNotFound)
	}
	return nil
}