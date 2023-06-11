package util

import (
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"time"
)

const (
	ATokenExpiredDuration = 2 * time.Hour
	RTokenExpiredDuration = 30 * 24 * time.Hour
)

type JWToken struct {
	SigningKey []byte
}

func NewJWT(sign string) *JWToken {
	return &JWToken{
		SigningKey: []byte(sign),
	}
}

// Generate 颁发token access token 和 refresh token
func (t *JWToken) Generate(claims model.Claims) (atoken, rtoken string, err error) {
	sc := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ATokenExpiredDuration).Unix(),
		Issuer:    "gogo",
		NotBefore: time.Now().Unix(), // 签名生效时间
	}
	claims.StandardClaims = sc
	atoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.SigningKey)

	// refresh token 不需要保存任何用户信息
	sc.ExpiresAt = time.Now().Add(RTokenExpiredDuration).Unix()
	rtoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, sc).SignedString(t.SigningKey)
	return
}

func (t *JWToken) keyfunc(token *jwt.Token) (any, error) {
	return t.SigningKey, nil
}

// Verify 验证Token
func (t *JWToken) Verify(tokenString string) (*model.Claims, error) {
	claims := new(model.Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, t.keyfunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token 验证不通过")
	}
	return claims, nil
}

// Refresh 通过 refresh token 刷新 atoken
func (t *JWToken) Refresh(atoken, rtoken string) (newAtoken, newRtoken, key string, err error) {
	// rtoken 无效直接返回
	if _, err = jwt.Parse(rtoken, t.keyfunc); err != nil {
		return
	}
	// 从旧access token 中解析出claims数据
	claim := new(model.Claims)
	_, err = jwt.ParseWithClaims(atoken, claim, t.keyfunc)
	// 判断错误是不是因为access token 正常过期导致的
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		at, rt, err := t.Generate(*claim)
		return at, rt, key, err
	}
	return
}
