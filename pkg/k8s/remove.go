package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

			fmt.Println(fmt.Sprintf("%s %s removida!", title, serviceName))
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

			fmt.Println(fmt.Sprintf("%s %s removida!", title, ingressName))
			return nil
		}
	}

	return nil
}
