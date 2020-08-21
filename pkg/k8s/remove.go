package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

var (
	functionResource = schema.GroupVersionResource{
		Group:    "openfaas.com",
		Version:  "v1",
		Resource: "functions",
	}

	kongPluginResource = schema.GroupVersionResource{
		Group:    "configuration.konghq.com",
		Version:  "v1",
		Resource: "kongplugins",
	}
)

func RemoveDeployment(client *kubernetes.Clientset, deployName, namespace, title string, verboseFlag bool) error {
	deployments := client.AppsV1().Deployments(namespace)
	listDeployments, _ := deployments.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações dos Deployments no cluster")
	}

	for _, deploy := range listDeployments.Items {
		if deploy.Name == deployName {
			err := deployments.Delete(context.TODO(), deployName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s %s removida!", title, deployName))
			return nil
		}
	}

	return nil
}

func RemoveService(client *kubernetes.Clientset, serviceName, namespace, title string, verboseFlag bool) error {
	services := client.CoreV1().Services(namespace)
	listServices, _ := services.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações dos Services no cluster")
	}

	for _, service := range listServices.Items {
		if service.Name == serviceName {
			err := services.Delete(context.TODO(), serviceName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			//fmt.Println(fmt.Sprintf("%s %s removida!", title, serviceName))
			return nil
		}
	}

	return nil
}

func RemoveIngress(client *kubernetes.Clientset, ingressName, namespace, title string, verboseFlag bool) error {
	ingresses := client.ExtensionsV1beta1().Ingresses(namespace)
	listIngresses, _ := ingresses.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações dos Ingresses no cluster")
	}

	for _, ingress := range listIngresses.Items {
		if ingress.Name == ingressName {
			err := ingresses.Delete(context.TODO(), ingressName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			//fmt.Println(fmt.Sprintf("%s %s removida!", title, ingressName))
			return nil
		}
	}

	return nil
}

func RemoveFunction(functionName, namespace, title, kubeconfig string, verboseFlag bool) error {
	client, err := GetDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	functions := client.Resource(functionResource)
	functionsList, _ := functions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações dos Functions no cluster")
	}

	for _, function := range functionsList.Items {
		if function.GetName() == functionName {
			err := functions.Namespace(namespace).Delete(context.TODO(), functionName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s %s removida!", title, functionName))
			return nil
		}
	}

	return nil
}

func RemovePlugin(kongPluginName, kubeconfig string, verboseFlag bool) error {
	client, err := GetDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	kongPlugin := client.Resource(kongPluginResource)
	kongPluginList, _ := kongPlugin.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações dos Kong Plugins no cluster")
	}

	for _, function := range kongPluginList.Items {
		if function.GetName() == kongPluginName {
			err := kongPlugin.Namespace("default").Delete(context.TODO(), kongPluginName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("Plugin %s removido!", kongPluginName))
			return nil
		}
	}

	return nil
}
