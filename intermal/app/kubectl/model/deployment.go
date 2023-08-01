package model

import corev1 "k8s.io/api/core/v1"

type Deployment struct {
	NamespaceWithName

	CPUMax    string `gorm:"NOT_NULL;COMMENT:'cpu最大值'" json:"cpuMax"`   // CPU, in cores. (500m = .5 cores)
	CPUMin    string `gorm:"NOT_NULL;COMMENT:'cpu最小'" json:"cpuMin"`    // CPU, in cores. (500m = .5 cores)
	MemoryMax string `gorm:"NOT_NULL;COMMENT:'内存最大值'" json:"memoryMax"` // Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	MemoryMin string `gorm:"NOT_NULL;COMMENT:'内存最小'" json:"memoryMin"`  // Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)

	Replicas int32 `gorm:"NOT_NULL;COMMENT:'副本数'" json:"replicas"`
	// PodPorts pod 开发端口
	Ports []PodPort `json:"ports"`
	// PodEnv  pod 使用的环境变量
	Env []PodEnv `json:"env"`

	// PullPolicy
	// - Always 总是拉取
	// - Never 只用本地，从不拉取
	// - IfNotPresent 默认值，本地有则使用本地镜像，否则拉取
	PullPolicy corev1.PullPolicy `gorm:"COMMENT:'镜像拉取策略'" json:"pullPolicy"`

	// RestartPolicy
	// - Always 当容器失效时，由kubelet自动重启该容器
	// - OnFailure 当容器终止运行且退出码不为0时。由kubelet自动重启该容器
	// - Never 不论容器状态如何，kubelet都不会重启该容器
	RestartPolicy corev1.RestartPolicy `gorm:"COMMENT:'Pod重启策略'" json:"restartPolicy"`

	// Recreate, Custom, Rolling
	ReleasePolicy string `gorm:"NOT_NULL;COMMENT:'发布策略'" json:"releasePolicy"`
	Image         string `gorm:"NOT_NULL;COMMENT:'镜像名称'" json:"image"`
}