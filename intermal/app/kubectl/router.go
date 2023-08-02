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

	// kubectl ingress
	ingressHandler := handler.NewIngressHandler(service.NewIngressService(k8sClient))
	//
	router.POST("kubectl/ingress/create", ingressHandler.Create)
	router.DELETE("kubectl/ingress/delete", ingressHandler.Delete)
	router.PUT("kubectl/ingress/update", ingressHandler.Update)
	router.POST("kubectl/ingress/list", ingressHandler.List)
	router.POST("kubectl/ingress/details", ingressHandler.Get)

	// kubectl volume
	volumeHandler := handler.NewVolumeHandler(service.NewVolumeService(k8sClient))
	//
	router.POST("kubectl/volume/create", volumeHandler.Create)
	router.DELETE("kubectl/volume/delete", volumeHandler.Delete)
	router.PUT("kubectl/volume/update", volumeHandler.Update)
	router.POST("kubectl/volume/list", volumeHandler.List)
	router.POST("kubectl/volume/details", volumeHandler.Get)

	// kubectl StatefulSet
	statefulSetHandler := handler.NewStatefulSetHandler(service.NewStatefulSetService(k8sClient))
	//
	router.POST("kubectl/statefulSet/create", statefulSetHandler.Create)
	router.DELETE("kubectl/statefulSet/delete", statefulSetHandler.Delete)
	router.PUT("kubectl/statefulSet/update", statefulSetHandler.Update)
	router.POST("kubectl/statefulSet/list", statefulSetHandler.List)
	router.POST("kubectl/statefulSet/details", statefulSetHandler.Get)
}