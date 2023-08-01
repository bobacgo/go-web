package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/app/admin/model/enum"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/mojocn/base64Captcha"
	"github.com/wenlng/go-captcha/captcha"
	"gorm.io/gorm"
	"time"
)

type ISystemService interface {
	Login(q model.LoginReq) (*model.LoginRsp, *g.Error)
	Refresh(q model.RefreshTokenVo) (model.RefreshTokenVo, error)
	Logout(username string) error
	Captcha() (model.CaptchaResponse, error)
	CaptchaV2() (model.CaptchaV2Response, *g.Error)
}

func NewSystemService(db *gorm.DB, userSvc IUserService, menuSvc IMenuService) ISystemService {
	return &systemService{
		db:          db,
		userService: userSvc,
		menuService: menuSvc,
	}
}

type systemService struct {
	db          *gorm.DB
	userService IUserService
	menuService IMenuService
}

var ErrUserDisable = errors.New("用户被禁用")

func (svc *systemService) Login(q model.LoginReq) (*model.LoginRsp, *g.Error) {

	// 要求
	// 1.校验验证码
	// 2.是否禁用（黑名单用户）
	// 3.校验密码
	// 4.登录次数（防止被暴力破解），登录次数超过5次直接禁用
	// 5.重置token（同时只能一个地方登录），如果有

	if !captchaStore.Verify(q.CaptchaKey, q.CaptchaCode, true) {
		return nil, g.NewError("验证码不正确")
	}
	user, gErr := svc.userService.FindWithRoleByUsername(svc.db, q.Username)
	if gErr != nil {
		return nil, g.NewError("用户名或密码错误")
	}
	if user.Status == enum.UserStatusDisable {
		return nil, g.WrapError(ErrUserDisable, "账号处于封禁状态")
	}
	if !passwordHelper.bcryptVerify(user.ID, q.Password, user.Password, func() {
		// 错误处理 达到一定错误次数后，禁用用户
		if err := svc.userService.UpdateStatus(model.UserUpdateStatusReq{
			IdReq: r.IdReq{ID: user.ID}, Status: enum.UserStatusDisable,
		}); err != nil {
			logger.Errorf("disable user [%s] err: %v", user.ID, err)
		}
	}) {
		return nil, g.NewError("用户名或密码错误")
	}

	at, rt, err := jwtService.Generate(user.SysUser)
	if err != nil {
		return nil, g.WrapError(err, "登录出错")
	}

	tree, gErr := svc.menuService.TreeByRole(user.RoleID)
	rsp := model.LoginRsp{
		User:   user,
		AToken: at,
		RToken: rt,
		Menus:  tree,
	}
	return &rsp, gErr
}

func (svc *systemService) Refresh(q model.RefreshTokenVo) (model.RefreshTokenVo, error) {
	at, rt, err := jwtService.Refresh(q.AToken, q.RToken)
	vo := model.RefreshTokenVo{
		AToken: at,
		RToken: rt,
	}
	return vo, err
}

func (svc *systemService) Logout(username string) error {
	err := jwtService.Remove(username)
	return err
}

func (svc *systemService) Captcha() (rsp model.CaptchaResponse, err error) {
	capt := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, captchaStore)
	rsp.CaptchaKey, rsp.ImgBase64, err = capt.Generate()
	return
}

func (svc *systemService) CaptchaV2() (model.CaptchaV2Response, *g.Error) {
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