package k8s

import (
	"fmt"

	"gopkg.in/gookit/color.v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("%s Error detected to kubeconfig, check: %s", color.Red.Text("[!]"), err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%s Error to create client, check kubeconfig: %s", color.Red.Text("[!]"), err.Error())
	}

	return client, nil
}

func GetDynamicClient(kubeconfig string) (dynamic.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("%s Error detected to kubeconfig, check: %s", color.Red.Text("[!]"), err.Error())
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%s Error to create client, check kubeconfig: %s", color.Red.Text("[!]"), err.Error())
	}

	return client, nil
}
