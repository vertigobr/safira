package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfig string) (*kubernetes.Clientset, error){
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error ao criar client, verifique o kubeconfig: %s", err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error ao verificar kubeconfig: %s", err.Error())
	}

	return client, nil
}
