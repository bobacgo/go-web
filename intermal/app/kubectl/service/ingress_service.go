package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model/enum"
	"github.com/gogoclouds/gogo/logger"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IIngressService interface {
	Create(context.Context, model.Ingress) error
	Delete(context.Context, model.NamespaceWithName) error
	Update(context.Context, model.Ingress) error
	Get(context.Context, model.NamespaceWithName) (*networkv1.Ingress, error)
	List(context.Context, string) ([]networkv1.Ingress, error)
}

func NewIngressService(k8sClient *kubernetes.Clientset) IIngressService {
	return &ingressService{
		k8sClient: k8sClient,
	}
}

type ingressService struct {
	k8sClient *kubernetes.Clientset
}

func (svc *ingressService) Create(ctx context.Context, o model.Ingress) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err == nil {
		return fmt.Errorf("ingress %s already exists", o.Name)
	} else {
		logger.Infow("get ingress:", o.Name, err.Error())
	}

	ingress := svc.makeIngress(o)
	_, err := svc.k8sClient.NetworkingV1().Ingresses(o.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
	return err
}

func (svc *ingressService) Delete(ctx context.Context, o model.NamespaceWithName) error {
	if _, err := svc.Get(ctx, o); err != nil {
		return fmt.Errorf("ingress %s not exists", o.Name)
	}

	err := svc.k8sClient.NetworkingV1().Ingresses(o.Namespace).Delete(ctx, o.Name, metav1.DeleteOptions{})
	return err
}

func (svc *ingressService) Update(ctx context.Context, o model.Ingress) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err != nil {
		return fmt.Errorf("ingress %s not exists", o.Name)
	}

	ingress := svc.makeIngress(o)
	_, err := svc.k8sClient.NetworkingV1().Ingresses(o.Namespace).Update(ctx, ingress, metav1.UpdateOptions{})
	return err
}

func (svc *ingressService) Get(ctx context.Context, o model.NamespaceWithName) (*networkv1.Ingress, error) {
	ingressInfo, err := svc.k8sClient.NetworkingV1().Ingresses(o.Namespace).Get(ctx, o.Name, metav1.GetOptions{})
	return ingressInfo, err
}

func (svc *ingressService) List(ctx context.Context, namespace string) ([]networkv1.Ingress, error) {
	res, err := svc.k8sClient.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []networkv1.Ingress{}, err
	}
	return res.Items, nil
}

/*
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:

	name: gogoclouds-ingress

spec:

		rules:
		- host: foo.bar.com
			http:
				paths:
				- pathType: Prefix
					path: "/"
					backend:
						service:
							name: service1
							port:
								number: 80
		- host: bar.foo.com
	    ...
*/
func (svc *ingressService) makeIngress(info model.Ingress) *networkv1.Ingress {
	ingress := new(networkv1.Ingress)
	ingress.TypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       enum.K8sKindIngress,
	}
	ingress.ObjectMeta = metav1.ObjectMeta{
		Name:      info.Name,
		Namespace: info.Namespace,
		Labels: map[string]string{
			"app-name": info.Name,
			"author":   "gogoclouds",
		},
	}
	// 使用 ingress-nginx
	className := "nginx"
	ingress.Spec = networkv1.IngressSpec{
		IngressClassName: &className,
		DefaultBackend:   nil, // 默认访问的服务
		TLS:              nil, // 开启https设置
		Rules:            svc.getIngressPath(info),
	}
	return ingress
}

func (svc *ingressService) getIngressPath(info model.Ingress) []networkv1.IngressRule {
	ingressRule := networkv1.IngressRule{
		Host: info.RouteHost,
	}
	ingressPaths := make([]networkv1.HTTPIngressPath, 0)
	for _, v := range info.RoutePath {
		pathType := networkv1.PathTypePrefix
		ingressPaths = append(ingressPaths, networkv1.HTTPIngressPath{
			Path:     v.PathName,
			PathType: &pathType,
			Backend: networkv1.IngressBackend{
				Service: &networkv1.IngressServiceBackend{
					Name: v.BackendService,
					Port: networkv1.ServiceBackendPort{
						Number: v.BackendServicePort,
					},
				},
			},
		})
	}
	ingressRule.IngressRuleValue = networkv1.IngressRuleValue{
		HTTP: &networkv1.HTTPIngressRuleValue{
			Paths: ingressPaths,
		},
	}
	return []networkv1.IngressRule{ingressRule}
}