package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin/service"
	"github.com/gogoclouds/go-web/intermal/common"
	"github.com/gogoclouds/go-web/intermal/util"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			reply.FailCode(ctx, r.TokenMission)
			ctx.Abort()
			return
		}
		jwtUtil := util.NewJWT(g.Conf.AppServiceKV()["authenticationKey"].(string))
		claims, err := jwtUtil.Verify(token)
		if err != nil {
			reply.FailCode(ctx, r.TokenInvalid)
			ctx.Abort()
			return
		}
		_, err = new(service.JwtService).Get(claims.Username)
		if err != nil {
			reply.FailCode(ctx, r.TokenInvalid)
			ctx.Abort()
			return
		}
		ctx.Set(common.GinContextClaimsKey, claims)
		ctx.Next()
	}
}
