package enum

type K8sKind string

const (
	K8sKindIngress               = "Ingress"
	K8sKindService               = "Service"
	K8sKindDeployment            = "Deployment"
	K8sKindStatefulSet           = "StatefulSet"
	K8sKindPersistentVolumeClaim = "PersistentVolumeClaim"
)