package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/admin"
	"github.com/gogoclouds/go-web/intermal/app/kubectl"
	"github.com/gogoclouds/go-web/intermal/middleware"
	"github.com/gogoclouds/gogo/g"
)

func loadRouter(e *gin.Engine) {
	e.MaxMultipartMemory = 300 << 20 //MB

	noAuthRouterGroup := e.Group("")
	admin.NoAuthRouterRegister(noAuthRouterGroup, g.DB)

	authRouterGroup := e.Group("")
	authRouterGroup.Use(middleware.JWTAuth())

	admin.RouterRegister(authRouterGroup, g.DB)

	k8sClient := getK8sClient()
	kubectl.RouterRegister(authRouterGroup, k8sClient)
}