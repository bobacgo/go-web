package model

import (
	"github.com/wenlng/go-captcha/captcha"
)

// LoginReq 登录请求
// 1.username + password
// 2.phone + password
// 3.phone + smsCode
// 4.第三方授权认证 TODO
type LoginReq struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaKey  string `json:"captchaKey" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

type LoginRsp struct {
	User   UserWithRole `json:"user"`
	AToken string       `json:"aToken"`
	RToken string       `json:"rToken"`
	Menus  []*SysMenu   `json:"menus"`
}

type RefreshTokenVo struct {
	AToken string `json:"aToken"`
	RToken string `json:"rToken"`
}

// =========================================
// Captcha
// =========================================

type CaptchaResponse struct {
	ImgBase64  string `json:"imgBase64"`
	CaptchaKey string `json:"captchaKey"`
}

type CaptchaV2Response struct {
	ImgBase64      string                  `json:"imgBase64"`
	ThumbImgBase64 string                  `json:"thumbImgBase64"`
	CaptchaKey     string                  `json:"captchaKey"`
	CharDots       map[int]captcha.CharDot `json:"charDots"`
}
