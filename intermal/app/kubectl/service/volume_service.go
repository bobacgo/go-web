package service

import (
	"context"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model"
	"github.com/gogoclouds/go-web/intermal/app/kubectl/model/enum"
	"github.com/gogoclouds/gogo/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IVolumeService interface {
	Create(context.Context, model.Volume) error
	Delete(context.Context, model.NamespaceWithName) error
	Update(context.Context, model.Volume) error
	Get(context.Context, model.NamespaceWithName) (*corev1.PersistentVolumeClaim, error)
	List(context.Context, string) ([]corev1.PersistentVolumeClaim, error)
}

func NewVolumeService(k8sClient *kubernetes.Clientset) IVolumeService {
	return &volumeService{
		k8sClient: k8sClient,
	}
}

type volumeService struct {
	k8sClient *kubernetes.Clientset
}

func (svc *volumeService) Create(ctx context.Context, o model.Volume) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err == nil {
		return fmt.Errorf("Volume %s already exists", o.Name)
	} else {
		logger.Infow("get Volume:", o.Name, err.Error())
	}

	pvc := svc.makeVolume(o)
	_, err := svc.k8sClient.CoreV1().PersistentVolumeClaims(o.Namespace).Create(ctx, pvc, metav1.CreateOptions{})
	return err
}

func (svc *volumeService) Delete(ctx context.Context, o model.NamespaceWithName) error {
	if _, err := svc.Get(ctx, o); err != nil {
		return fmt.Errorf("volume %s not exists", o.Name)
	}

	err := svc.k8sClient.CoreV1().PersistentVolumeClaims(o.Namespace).Delete(ctx, o.Name, metav1.DeleteOptions{})
	return err
}

func (svc *volumeService) Update(ctx context.Context, o model.Volume) error {
	if _, err := svc.Get(ctx, o.NamespaceWithName); err != nil {
		return fmt.Errorf("volume %s not exists", o.Name)
	}

	pvc := svc.makeVolume(o)
	_, err := svc.k8sClient.CoreV1().PersistentVolumeClaims(o.Namespace).Update(ctx, pvc, metav1.UpdateOptions{})
	return err
}

func (svc *volumeService) Get(ctx context.Context, o model.NamespaceWithName) (*corev1.PersistentVolumeClaim, error) {
	pvcInfo, err := svc.k8sClient.CoreV1().PersistentVolumeClaims(o.Namespace).Get(ctx, o.Name, metav1.GetOptions{})
	return pvcInfo, err
}

func (svc *volumeService) List(ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	res, err := svc.k8sClient.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []corev1.PersistentVolumeClaim{}, err
	}
	return res.Items, err
}

/*
apiVersion: v1
kind: PersistentVolume
metadata:

	name: local-storage-pv-1
	namespace: elasticsearch
	labels:
	  name: local-storage-pv-1

spec:

	capacity:
	  storage: 1Gi
	accessModes:
	- ReadWriteOnce
	persistentVolumeReclaimPolicy: Retain
	storageClassName: local-storage
	local:
	  path: /data/es
	nodeAffinity:
	  required:
	    nodeSelectorTerms:
	    - matchExpressions:
	      - key: kubernetes.io/hostname
	        operator: In
	        values:
	        - master1
*/
func (svc *volumeService) makeVolume(info model.Volume) *corev1.PersistentVolumeClaim {
	pvc := new(corev1.PersistentVolumeClaim)
	pvc.TypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       enum.K8sKindPersistentVolumeClaim,
	}
	pvc.ObjectMeta = metav1.ObjectMeta{
		Name:      info.Name,
		Namespace: info.Namespace,
		Labels:    nil,
		Annotations: map[string]string{
			"pv.kubernetes.io/bound-by-controller":          "yes",
			"volume.beta.kubernetes.io/storage-provisioner": "rbd.csi.ceph.com",
		},
	}
	pvc.Spec = corev1.PersistentVolumeClaimSpec{
		AccessModes:      info.AccessMode,
		Resources:        svc.getResources(info),
		StorageClassName: &info.StorageClass,
		VolumeMode:       &info.VolumeMode,
	}
	return pvc
}

// 获取资源配置
func (svc *volumeService) getResources(info model.Volume) (res corev1.ResourceRequirements) {
	res.Requests = corev1.ResourceList{
		corev1.ResourceStorage: resource.MustParse(info.Storage),
	}
	return
}