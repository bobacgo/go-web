package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model/enum"
	"github.com/gogoclouds/gogo/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IDeploymentService interface {
	Create(context.Context, model.Deployment) error
	Delete(context.Context, model.NamespaceWithName) error
	Update(context.Context, model.Deployment) error
	Get(context.Context, model.NamespaceWithName) (*appsv1.Deployment, error)
	List(context.Context, string) ([]appsv1.Deployment, error)
}

func NewDeploymentService(k8sClient *kubernetes.Clientset) IDeploymentService {
	return &deploymentService{
		k8sClient: k8sClient,
	}
}

type deploymentService struct {
	k8sClient *kubernetes.Clientset
}

func (svc *deploymentService) Create(ctx context.Context, o model.Deployment) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err == nil {
		return fmt.Errorf("deployment %s already exists", o.Name)
	} else {
		logger.Infow("get deployment:", o.Name, err.Error())
	}

	deployment := svc.makeDeployment(o)
	_, err := svc.k8sClient.AppsV1().Deployments(o.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

func (svc *deploymentService) Delete(ctx context.Context, req model.NamespaceWithName) error {
	if _, err := svc.Get(ctx, req); err != nil {
		return fmt.Errorf("deployment %s not exists", req.Name)
	}
	err := svc.k8sClient.AppsV1().Deployments(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	return err
}

func (svc *deploymentService) Update(ctx context.Context, o model.Deployment) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err != nil {
		return fmt.Errorf("deployment %s not exists", o.Name)
	}
	deployment := svc.makeDeployment(o)
	_, err := svc.k8sClient.AppsV1().Deployments(o.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

func (svc *deploymentService) Get(ctx context.Context, o model.NamespaceWithName) (*appsv1.Deployment, error) {
	res, err := svc.k8sClient.AppsV1().Deployments(o.Namespace).Get(ctx, o.Name, metav1.GetOptions{})
	return res, err
}

func (svc *deploymentService) List(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	res, err := svc.k8sClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []appsv1.Deployment{}, err
	}
	return res.Items, nil
}

/*
apiVersion: apps/v1
kind: Deployment
metadata:

	namespace: gogo
	name: go-web
	labels:
	  app: go-web
	  version: v1

spec:

	replicas: 1
	selector:
	  matchLabels:
	    app: go-web
	template:
	  metadata:
	    labels:
	      app: go-web
	  spec:
	    containers:
	      - name: go
	        image: xxx
	        ports:
	          - containerPort: 8888
*/
func (svc *deploymentService) makeDeployment(info model.Deployment) *appsv1.Deployment {
	deploy := new(appsv1.Deployment)
	deploy.TypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       enum.K8sKindDeployment,
	}
	deploy.ObjectMeta = metav1.ObjectMeta{
		Name:      info.Name,
		Namespace: info.Namespace,
		Labels: map[string]string{
			"app-name": info.Name,
			"author":   "gogoclouds",
		},
	}
	deploy.Name = info.Name
	deploy.Spec = appsv1.DeploymentSpec{
		Replicas: &info.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": info.Name,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app-name": info.Name,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            info.Name,
						Image:           info.Image,
						Ports:           svc.getContainerPorts(info.Ports),
						Env:             svc.getEnv(info.Env),
						Resources:       svc.getResources(info),
						ImagePullPolicy: info.PullPolicy,
					},
				},
			},
		},
	}
	return deploy
}

func (svc *deploymentService) getResources(info model.Deployment) (res corev1.ResourceRequirements) {
	// 最大使用资源
	res.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(info.CPUMax),
		corev1.ResourceMemory: resource.MustParse(info.MemoryMax),
	}
	// 最少使用资源
	res.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(info.CPUMin),
		corev1.ResourceMemory: resource.MustParse(info.MemoryMin),
	}
	return
}

func (svc *deploymentService) getEnv(podEnv []model.PodEnv) (env []corev1.EnvVar) {
	for _, v := range podEnv {
		env = append(env, corev1.EnvVar{
			Name:  v.Key,
			Value: v.Value,
		})
	}
	return
}

func (svc *deploymentService) getContainerPorts(podPorts []model.PodPort) (ports []corev1.ContainerPort) {
	for _, v := range podPorts {
		ports = append(ports, corev1.ContainerPort{
			Name:          fmt.Sprintf("port-%d", v.ContainerPort),
			ContainerPort: v.ContainerPort,
			Protocol:      v.Protocol,
		})
	}
	return
}