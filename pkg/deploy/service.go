// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"context"
	"fmt"
	"github.com/vertigobr/safira/pkg/k8s"
	"gopkg.in/gookit/color.v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
)

type serviceSpec struct {
	Type         string            `yaml:"type,omitempty"`
	Selector     map[string]string `yaml:"selector,omitempty"`
	ExternalName string            `yaml:"externalName,omitempty"`
	Ports        []servicePort     `yaml:"ports,omitempty"`
}

type servicePort struct {
	Protocol   string `yaml:"protocol,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	TargetPort int    `yaml:"targetPort,omitempty"`
}

func (k *K8sYaml) MountService(serviceName, hostname, env string, isFunction bool) error {
	stack, err := s.LoadStackFile(env)
	if err != nil {
		return err
	}

	var port int
	if len(hostname) > 1 {
		_, port, err = getGatewayPort(hostname)
	} else {
		_, port, err = getGatewayPort(stack.Hostname)
	}

	if err != nil {
		return err
	}

	spec := getServiceSpec(serviceName, port, isFunction)
	annotations, err := getServiceAnnotations(serviceName, stack.Functions, isFunction)
	if err != nil {
		return err
	}

	*k = K8sYaml{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: metadata{
			Name: serviceName,
			Labels: map[string]string{
				"app": serviceName,
			},
			Annotations: annotations,
		},
		Spec: spec,
	}

	return nil
}

func getServiceSpec(serviceName string, port int, isFunction bool) (spec serviceSpec) {
	if isFunction {
		spec = serviceSpec{
			Type:         "ExternalName",
			ExternalName: "gateway",
			Ports: []servicePort{
				{
					Port: port,
				},
			},
		}
	} else {
		spec = serviceSpec{
			Type: "NodePort",
			Selector: map[string]string{
				"app": serviceName,
			},
			Ports: []servicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: 8080,
				},
			},
		}
	}

	return
}

func getServiceAnnotations(serviceName string, functions map[string]s.Function, isFunction bool) (map[string]string, error) {
	annotations := make(map[string]string)

	if isFunction {
		for functionName, function := range functions {
			if functionName == serviceName {
				for pluginName, plugin := range function.Plugins {
					if len(plugin.Type) == 0 || plugin.Type == "service" {
						annotations["konghq.com/plugins"] = fmt.Sprintf("%s-%s", functionName, pluginName)
					}
				}
			}
		}
	}

	repoName, err := utils.GetCurrentFolder()
	if err != nil {
		return nil, err
	}

	annotations["safira.io/repository"] = repoName
	annotations["safira.io/function"] = serviceName

	return annotations, nil
}

func AddPluginInAnnotationsService(name, namespace, pluginName, kubeconfig string, verboseFlag bool) error {
	client, err := k8s.GetClient(kubeconfig)
	if err != nil {
		if verboseFlag {
			fmt.Println(err.Error())
		}

		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	services := client.CoreV1().Services(namespace)
	service, err := services.Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("")
	}

	service.ObjectMeta.Annotations["konghq.com/plugins"] = pluginName
	_, err = services.Update(context.TODO(), service, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("")
	}

	return nil
}
