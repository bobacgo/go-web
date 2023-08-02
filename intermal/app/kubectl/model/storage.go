package model

import corev1 "k8s.io/api/core/v1"

type Storage struct {
	//存储名称
	Name string `json:"name"`
	//存储的大小
	Size string `json:"size"`
	//存储需要挂载的目录
	MountPath string `json:"path"`
	//存储创建的类型
	StorageClass string `json:"storageClass"`
	//存储的权限
	AccessMode []corev1.PersistentVolumeAccessMode `json:"accessMode"`
}