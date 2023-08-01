package model

import corev1 "k8s.io/api/core/v1"

type Service struct {
	NamespaceWithName
	PodName string `gorm:"NOT_NULL;COMMENT:'绑定Pod名称'" json:"podName"`
	// ServiceType
	// - ClusterIP 只能通过集群IP在集群内部访问
	// - NodePort 对集群外暴露端口方式，给集群外服务访问
	// - LoadBalancer 需要云提供商支持
	// - ExternalName 包含引用外部服务
	ServiceType             corev1.ServiceType `gorm:"NOT_NULL" json:"serviceType"`
	ServiceTypeExternalName string             `gorm:"COMMENT:'service类型为ExternalName时'" json:"serviceTypeExternal"`
	// ServicePort service 上设置端口
	Ports []ServicePort `json:"ports"`
}