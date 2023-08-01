package kubectl

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/handler"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/service"
	"k8s.io/client-go/kubernetes"
)

func RouterRegister(router *gin.RouterGroup, k8sClient *kubernetes.Clientset) {
	// kubectl pod
	deploymentHandler := handler.NewDeploymentHandler(service.NewDeploymentService(k8sClient))
	//
	router.POST("kubectl/deployment/create", deploymentHandler.Create)
	router.DELETE("kubectl/deployment/delete", deploymentHandler.Delete)
	router.PUT("kubectl/deployment/update", deploymentHandler.Update)
	router.POST("kubectl/deployment/list", deploymentHandler.List)
	router.POST("kubectl/deployment/details", deploymentHandler.Get)

	// kubectl service
	serviceHandler := handler.NewServiceHandler(service.NewSvcService(k8sClient))
	//
	router.POST("kubectl/service/create", serviceHandler.Create)
	router.DELETE("kubectl/service/delete", serviceHandler.Delete)
	router.PUT("kubectl/service/update", serviceHandler.Update)
	router.POST("kubectl/service/list", serviceHandler.List)
	router.POST("kubectl/service/details", serviceHandler.Get)
}