package model

import corev1 "k8s.io/api/core/v1"

type ServicePort struct {
	Port       int32 `gorm:"COMMENT:'service的端口'" json:"port"`
	TargetPort int   `gorm:"COMMENT:'pod 中需要映射的port地址'" json:"targetPort"`
	NodePort   int32 `gorm:"COMMENT:'NodePort的模式下进行设置'" json:"nodePort"`
	// Protocol TCP|UDP|SCTP
	PortProtocol corev1.Protocol `json:"portProtocol"`
}