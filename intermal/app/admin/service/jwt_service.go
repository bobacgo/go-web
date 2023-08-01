package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/go-web/pkg/util"
	"github.com/gogoclouds/gogo/g"
	"time"
)

type IJwtService interface {
	Generate(user model.SysUser) (string, string, error)
	Refresh(aToken, rToken string) (string, string, error)
	Set(username, token string) error
	Get(username string) (string, error)
	Remove(username string) error
}

var jwtService IJwtService = new(JwtService)

type JwtService struct{}

func (svc *JwtService) Generate(user model.SysUser) (string, string, error) {
	signKey := g.Conf.AppServiceKV()["authenticationKey"]
	claims := model.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		RoleID:   user.RoleID,
	}
	aToken, rToken, err := util.NewJWT(signKey.(string)).Generate(claims)
	if err != nil {
		return "", "", err
	}
	err = jwtService.Set(user.Username, aToken)
	return aToken, rToken, err
}

func (svc *JwtService) Refresh(aToken, rToken string) (string, string, error) {
	signKey := g.Conf.AppServiceKV()["authenticationKey"]
	atoken, rtoken, key, err := util.NewJWT(signKey.(string)).Refresh(aToken, rToken)
	if err == nil {
		if err := svc.Set(key, atoken); err != nil {
			return "", "", err
		}
	}
	return atoken, rtoken, err
}

func (svc *JwtService) Set(username, token string) error {
	err := g.CacheDB.Set(context.Background(), svc.key(username), token, 2*time.Hour).Err()
	return err
}

func (svc *JwtService) Get(username string) (string, error) {
	result, err := g.CacheDB.Get(context.Background(), svc.key(username)).Result()
	return result, err
}

func (svc *JwtService) Remove(username string) error {
	err := g.CacheDB.Del(context.Background(), svc.key(username)).Err()
	return err
}

func (svc *JwtService) key(username string) string {
	return fmt.Sprintf(common.RedisKeyJwt, username)
}