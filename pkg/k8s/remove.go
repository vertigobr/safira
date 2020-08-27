package k8s

import (
	"context"
	"fmt"

	"gopkg.in/gookit/color.v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func RemoveDeployment(name, namespace, kubeconfig string, verboseFlag bool) error {
	client, err := GetClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	deployments := client.AppsV1().Deployments(namespace)
	listDeployments, _ := deployments.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting deployments info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, deploy := range listDeployments.Items {
		if deploy.Name == name {
			err := deployments.Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("%s Deployment %s removed\n", color.Green.Text("[+]"), name)
			return nil
		}
	}

	return nil
}

func RemoveService(name, namespace, kubeconfig string, verboseFlag bool) error {
	client, err := GetClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	services := client.CoreV1().Services(namespace)
	listServices, _ := services.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting services info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, service := range listServices.Items {
		if service.Name == name {
			err := services.Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			if verboseFlag {
				fmt.Printf("%s Service %s removed\n", color.Blue.Text("[v]"), name)
			}

			return nil
		}
	}

	return nil
}

func RemoveIngress(name, namespace, kubeconfig string, verboseFlag bool) error {
	client, err := GetClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	ingresses := client.ExtensionsV1beta1().Ingresses(namespace)
	listIngresses, _ := ingresses.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting ingresses info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, ingress := range listIngresses.Items {
		if ingress.Name == name {
			err := ingresses.Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			if verboseFlag {
				fmt.Printf("%s Ingress %s removed\n", color.Blue.Text("[v]"), name)
			}

			return nil
		}
	}

	return nil
}

func RemoveConfigmap(name, namespace, kubeconfig string, verboseFlag bool) error {
	client, err := GetClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	configmaps := client.CoreV1().ConfigMaps(namespace)
	listConfigmap, _ := configmaps.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting configmaps info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, service := range listConfigmap.Items {
		if service.Name == name {
			err := configmaps.Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			if verboseFlag {
				fmt.Printf("%s Configmap %s removed\n", color.Blue.Text("[v]"), name)
			}

			return nil
		}
	}

	return nil
}

func RemoveFunction(name, namespace, kubeconfig string, verboseFlag bool) error {
	client, err := GetDynamicClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	functions := client.Resource(functionResource)
	functionsList, _ := functions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting functions info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, function := range functionsList.Items {
		if function.GetName() == name {
			err := functions.Namespace(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("%s Function %s removed\n", color.Green.Text("[+]"), name)
			return nil
		}
	}

	return nil
}

func RemoveKongPlugin(name, kubeconfig string, verboseFlag bool) error {
	client, err := GetDynamicClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	kongPlugin := client.Resource(kongPluginResource)
	kongPluginList, _ := kongPlugin.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Getting kongplugins info in the cluster\n", color.Blue.Text("[v]"))
	}

	for _, function := range kongPluginList.Items {
		if function.GetName() == name {
			err := kongPlugin.Namespace("default").Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("%s KongPlugin %s removed\n", color.Green.Text("[+]"), name)
			return nil
		}
	}

	return nil
}
