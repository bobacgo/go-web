package model

import corev1 "k8s.io/api/core/v1"

type PodPort struct {
	ContainerPort int32 `gorm:"NOT_NULL" json:"container_port"`
	// Protocol default TCP TCP|UDP|SCTP
	Protocol corev1.Protocol `json:"protocol"`
}