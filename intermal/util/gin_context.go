package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/go-web/intermal/common"
)

func ContextClaims(ctx *gin.Context) *model.Claims {
	v, ok := ctx.Get(common.GinContextClaimsKey)
	if ok {
		return nil
	}
	claims := v.(model.Claims)
	return &claims
}

func ContextUsername(ctx *gin.Context) string {
	claims := ContextClaims(ctx)
	if claims == nil {
		return ""
	}
	return claims.Username
}
