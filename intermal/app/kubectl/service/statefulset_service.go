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

type IStatefulSetService interface {
	Create(context.Context, model.StatefulSet) error
	Delete(context.Context, model.NamespaceWithName) error
	Update(context.Context, model.StatefulSet) error
	Get(context.Context, model.NamespaceWithName) (*appsv1.StatefulSet, error)
	List(context.Context, string) ([]appsv1.StatefulSet, error)
}

func NewStatefulSetService(k8sClient *kubernetes.Clientset) IStatefulSetService {
	return &statefulSetService{
		k8sClient: k8sClient,
	}
}

type statefulSetService struct {
	k8sClient *kubernetes.Clientset
}

func (svc *statefulSetService) Create(ctx context.Context, o model.StatefulSet) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err == nil {
		return fmt.Errorf("statefulSet %s already exists", o.Name)
	} else {
		logger.Infow("get statefulSet:", o.Name, err.Error())
	}

	ss := svc.makeStatefulSet(o)
	_, err := svc.k8sClient.AppsV1().StatefulSets(o.Namespace).Create(ctx, ss, metav1.CreateOptions{})
	return err
}

func (svc *statefulSetService) Delete(ctx context.Context, o model.NamespaceWithName) error {
	if _, err := svc.Get(ctx, o); err != nil {
		return fmt.Errorf("statefulSet %s not exists", o.Name)
	}

	err := svc.k8sClient.AppsV1().StatefulSets(o.Namespace).Delete(ctx, o.Name, metav1.DeleteOptions{})
	return err
}

func (svc *statefulSetService) Update(ctx context.Context, o model.StatefulSet) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err != nil {
		return fmt.Errorf("statefulSet %s not exists", o.Name)
	}

	service := svc.makeStatefulSet(o)
	_, err := svc.k8sClient.AppsV1().StatefulSets(o.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	return err
}

func (svc *statefulSetService) Get(ctx context.Context, o model.NamespaceWithName) (*appsv1.StatefulSet, error) {
	info, err := svc.k8sClient.AppsV1().StatefulSets(o.Namespace).Get(ctx, o.Name, metav1.GetOptions{})
	return info, err
}

func (svc *statefulSetService) List(ctx context.Context, namespace string) ([]appsv1.StatefulSet, error) {
	info, err := svc.k8sClient.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []appsv1.StatefulSet{}, err
	}
	return info.Items, err
}

// 根据info信息设置值
func (svc *statefulSetService) makeStatefulSet(info model.StatefulSet) *appsv1.StatefulSet {
	statefulSet := new(appsv1.StatefulSet)
	statefulSet.TypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       enum.K8sKindStatefulSet,
	}
	//设置详情
	statefulSet.ObjectMeta = metav1.ObjectMeta{
		Name:      info.Name,
		Namespace: info.Namespace,
		//设置label标签
		Labels: map[string]string{
			"app-name": info.Name,
			"author":   "gogoclouds",
		},
	}
	statefulSet.Name = info.Name
	var terminationGracePeriodSeconds int64 = 10
	statefulSet.Spec = appsv1.StatefulSetSpec{
		//副本数
		Replicas: &info.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": info.Name,
			},
		},
		//设置容器模版
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app-name": info.Name,
				},
			},
			//设置容器详情
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  info.Name,
						Image: info.Image,
						//获取容器的端口
						Ports: svc.getContainerPorts(info.Ports),
						//获取环境变量
						Env: svc.getEnv(info.PodEnv),
						//获取容器CPU，内存
						Resources: svc.getResources(info),
						//设置挂载目录
						VolumeMounts: svc.getMounts(info.Storage),
					},
				},
				//不能设置为0，这样不安全
				//https://kubernetes.io/docs/tasks/run-application/force-delete-stateful-set-pod/
				TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
				//私有仓库设置密钥
				ImagePullSecrets: nil,
			},
		},
		VolumeClaimTemplates: svc.getPVC(info),
		ServiceName:          info.Name,
	}
	return statefulSet
}

// 获取pvc
func (svc *statefulSetService) getPVC(info model.StatefulSet) (pvcAll []corev1.PersistentVolumeClaim) {
	if len(info.Storage) == 0 {
		return
	}
	for _, v := range info.Storage {
		pvc := &corev1.PersistentVolumeClaim{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       enum.K8sKindPersistentVolumeClaim,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      v.Name,
				Namespace: info.Namespace,
				Annotations: map[string]string{
					"pv.kubernetes.io/bound-by-controller":          "yes",
					"volume.beta.kubernetes.io/storage-provisioner": "rbd.csi.ceph.com",
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      v.AccessMode,
				Resources:        svc.getPvcResource(v.Size),
				VolumeName:       v.Name,
				StorageClassName: &v.StorageClass,
			},
		}
		pvcAll = append(pvcAll, *pvc)
	}
	return
}

func (svc *statefulSetService) getPvcResource(size string) (source corev1.ResourceRequirements) {
	source.Requests = corev1.ResourceList{
		corev1.ResourceStorage: resource.MustParse(size),
	}
	return
}

// 设置存储路径
func (svc *statefulSetService) getMounts(vs []model.Storage) (mount []corev1.VolumeMount) {
	if len(vs) == 0 {
		return
	}
	for _, v := range vs {
		mt := &corev1.VolumeMount{
			Name:      v.Name,
			MountPath: v.MountPath,
		}
		mount = append(mount, *mt)
	}
	return
}

func (svc *statefulSetService) getResources(info model.StatefulSet) (res corev1.ResourceRequirements) {
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

func (svc *statefulSetService) getEnv(podEnv []model.PodEnv) (env []corev1.EnvVar) {
	for _, v := range podEnv {
		env = append(env, corev1.EnvVar{
			Name:  v.Key,
			Value: v.Value,
		})
	}
	return
}

func (svc *statefulSetService) getContainerPorts(podPorts []model.PodPort) (ports []corev1.ContainerPort) {
	for _, v := range podPorts {
		ports = append(ports, corev1.ContainerPort{
			Name:          fmt.Sprintf("port-%d", v.ContainerPort),
			ContainerPort: v.ContainerPort,
			Protocol:      v.Protocol,
		})
	}
	return
}