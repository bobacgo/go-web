package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model/enum"
	"github.com/gogoclouds/gogo/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type ISvcService interface {
	Create(context.Context, model.Service) error
	Delete(context.Context, model.NamespaceWithName) error
	Update(context.Context, model.Service) error
	Get(context.Context, model.NamespaceWithName) (*corev1.Service, error)
	List(context.Context, string) ([]corev1.Service, error)
}

func NewSvcService(k8sClient *kubernetes.Clientset) ISvcService {
	return &svcService{
		k8sClient: k8sClient,
	}
}

type svcService struct {
	k8sClient *kubernetes.Clientset
}

func (svc *svcService) Create(ctx context.Context, o model.Service) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err == nil {
		return fmt.Errorf("service %s already exists", o.Name)
	} else {
		logger.Infow("get service:", o.Name, err.Error())
	}

	service := svc.makeService(o)
	_, err := svc.k8sClient.CoreV1().Services(o.Namespace).Create(ctx, service, metav1.CreateOptions{})
	return err

}

func (svc *svcService) Delete(ctx context.Context, o model.NamespaceWithName) error {
	if _, err := svc.Get(ctx, o); err != nil {
		return fmt.Errorf("service %s not exists", o.Name)
	}

	err := svc.k8sClient.CoreV1().Services(o.Namespace).Delete(ctx, o.Name, metav1.DeleteOptions{})
	return err
}

func (svc *svcService) Update(ctx context.Context, o model.Service) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err != nil {
		return fmt.Errorf("service %s not exists", o.Name)
	}

	service := svc.makeService(o)
	_, err := svc.k8sClient.CoreV1().Services(o.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	return err
}

func (svc *svcService) Get(ctx context.Context, req model.NamespaceWithName) (*corev1.Service, error) {
	serviceInfo, err := svc.k8sClient.CoreV1().Services(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	return serviceInfo, err
}

func (svc *svcService) List(ctx context.Context, namespace string) ([]corev1.Service, error) {
	res, err := svc.k8sClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []corev1.Service{}, err
	}
	return res.Items, nil
}

/*
apiVersion: v1
kind: Service
metadata:

	name: go-server
	namespace: gogo
	labels:
	  app: go-server

spec:

	type: NodePort
	selector:
	  app: go-web
	ports:
	  - port: 8888
	    targetPort: 8888
	    name: tcp
*/
func (svc *svcService) makeService(o model.Service) *corev1.Service {
	k8sService := new(corev1.Service)
	k8sService.TypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       enum.K8sKindService,
	}
	// 设置服务基础信息
	k8sService.ObjectMeta = metav1.ObjectMeta{
		Name:      o.Name,
		Namespace: o.Namespace,
		Labels: map[string]string{
			"app-name": o.Name,
			"author":   "gogoclouds",
		},
	}
	//设置服务spec信息
	k8sService.Spec = corev1.ServiceSpec{
		Ports: svc.getSvcPort(o.Ports),
		Selector: map[string]string{
			"app-name": o.PodName,
		},
		Type: o.ServiceType,
	}
	return k8sService
}

func (svc *svcService) getSvcPort(svcPort []model.ServicePort) (ports []corev1.ServicePort) {
	for _, v := range svcPort {
		ports = append(ports, corev1.ServicePort{
			Name:       fmt.Sprintf("port-%d", v.Port),
			Port:       v.Port,
			TargetPort: intstr.FromInt(v.TargetPort),
			Protocol:   v.PortProtocol,
		})
	}
	return
}