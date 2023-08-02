package model

import corev1 "k8s.io/api/core/v1"

type Volume struct {
	NamespaceWithName
	//AccessMode 存储的访问模式
	// PersistentVolumeAccessMode
	// - ReadWriteOnce 读写模式只能挂载一个节点
	// - ReadOnlyMany 只读模式挂载多个节点
	// - ReadWriteMany 读写模式挂载多个节点
	// - ReadWriteOncePod 读写模式仅挂载一个Pod
	AccessMode []corev1.PersistentVolumeAccessMode `json:"accessMode"`
	// sc 的class name
	StorageClass string `json:"storageClass"`
	// Storage 资源大小
	Storage string `json:"storage"`
	// VolumeMode 存储类型
	// - Block
	// - Filesystem
	VolumeMode corev1.PersistentVolumeMode `json:"volumeMode"`
}