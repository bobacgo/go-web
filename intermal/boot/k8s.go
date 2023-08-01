package boot

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func getK8sClient() *kubernetes.Clientset {
	//创建 config 实例
	// docker 内 C:/Users/mr/.kube/config:/root/.kube/config
	hdir := homedir.HomeDir()
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(hdir, ".kube", "config"))
	if err != nil {
		panic(err)
	}

	//在集群中使用
	//config , err := rest.InClusterConfig()
	//if err!=nil {
	//	panic(err)
	//}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err)
	}
	return clientset
}