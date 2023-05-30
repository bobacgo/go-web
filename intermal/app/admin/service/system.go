package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/logger"
	"github.com/mojocn/base64Captcha"
	"github.com/wenlng/go-captcha/captcha"
	"time"
)

type ISystem interface {
	Login(q model.LoginReq) ([]*model.SysMenu, *g.Error)
	Logout()
	Captcha() (model.CaptchaResponse, error)
	CaptchaV2() (model.CaptchaV2Response, *g.Error)
}

var menuService IMenuService = new(MenuService)

type System struct{}

var ErrUserDisable = errors.New("用户被禁用")

func (s *System) Login(q model.LoginReq) ([]*model.SysMenu, *g.Error) {

	// 要求
	// 1.校验验证码
	// 2.是否禁用（黑名单用户）
	// 3.校验密码
	// 4.登录次数（防止被暴力破解），登录次数超过5次直接禁用
	// 5.重置token（同时只能一个地方登录），如果有

	if !captchaStore.Verify(q.CaptchaKey, q.CaptchaCode, true) {
		return nil, g.NewError("验证码不正确")
	}
	user, gErr := userService.findByUsername(g.DB, q.Username)
	if gErr != nil {
		logger.Errorf("[%s] %v", q.Username, gErr.Error())
		return nil, g.NewError("用户名或密码错误")
	}
	if user.Status == enum.UserStatusDisable {
		return nil, g.WrapError(ErrUserDisable, "账号处于封禁状态")
	}
	//up := strings.SplitN(user.Password, "$", 1)
	tree, gErr := menuService.TreeByRole("a70e443047f0403094d406b5a3c78880")
	return tree, gErr
}

func (s *System) Logout() {
	//TODO implement me
	panic("implement me")
}

func (s *System) Captcha() (rsp model.CaptchaResponse, err error) {
	capt := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, captchaStore)
	rsp.CaptchaKey, rsp.ImgBase64, err = capt.Generate()
	return
}

func (s *System) CaptchaV2() (model.CaptchaV2Response, *g.Error) {
	// 生成验证码
	var (
		rsp model.CaptchaV2Response
		err error
	)
	rsp.CharDots, rsp.ImgBase64, rsp.ThumbImgBase64, rsp.CaptchaKey, err = captcha.GetCaptcha().Generate()
	return rsp, g.WrapError(err, "生成验证码出错")
}

// ===============================================
// Captcha Redis Store
// ===============================================

var captchaStore base64Captcha.Store = new(captchaRedisStore)

type captchaRedisStore struct{}

func (c *captchaRedisStore) Set(id string, value string) error {
	err := g.CacheDB.Set(context.Background(), c.key(id), value, time.Minute).Err()
	return err
}

func (c *captchaRedisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	result, err := g.CacheDB.Get(ctx, c.key(id)).Result()
	if err != nil {
		logger.Error(err)
	}
	if clear {
		if err = g.CacheDB.Del(ctx, c.key(id)).Err(); err != nil {
			logger.Error(err)
		}
	}
	return result
}

func (c *captchaRedisStore) Verify(id, answer string, clear bool) bool {
	captchaCode := c.Get(id, clear)
	return captchaCode == answer
}

func (c *captchaRedisStore) key(id string) string {
	return fmt.Sprintf(common.RedisKeyCaptchaFmt, id)
}