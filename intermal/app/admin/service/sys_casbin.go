package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/gogoclouds/go-web/intermal/app/admin/dao"
	"github.com/gogoclouds/gogo/logger"
	"sync"
)

type ICasbin interface {
	Casbin() *casbin.SyncedCachedEnforcer
}

var casbinDao dao.ICasbin = new(dao.Casbin)

type Casbin struct{}

var (
	once                 sync.Once
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
)

func (svc Casbin) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		var err error
		if syncedCachedEnforcer, err = casbinDao.Casbin(); err != nil {
			logger.Error(err)
			return
		}
	})
	return syncedCachedEnforcer
}

func (svc Casbin) RemoveFilteredPolicy(v int, p ...string) bool {
	e := svc.Casbin()
	ok, err := e.RemoveFilteredPolicy(v, p...)
	if err != nil {
		logger.Error(err)
	}
	return ok
}

func (svc Casbin) UpdateCasbinApi(menuID string, path, method string) error {
	//casbinDao.UpdateByMenuID(menuID, path, method)
	e := svc.Casbin()
	err := e.InvalidateCache()
	return err
}