// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"

	s "github.com/vertigobr/safira/pkg/stack"
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
	annotations := getServiceMetadata(serviceName, stack.Functions, isFunction)

	*k = K8sYaml{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: metadata{
			Name:   serviceName,
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
			Type: "ExternalName",
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
					Protocol: "TCP",
					Port: 80,
					TargetPort: 8080,
				},
			},
		}
	}

	return
}

func getServiceMetadata(serviceName string, functions map[string]s.Function, isFunction bool) map[string]string {
	annotations := make(map[string]string)

	if isFunction {
		for functionName, function := range functions {
			if functionName == serviceName {
				for pluginName, plugin := range function.Plugins {
					if len(plugin.Type) == 0 || plugin.Type == "service" {
						annotations["plugins.konghq.com"] = fmt.Sprintf("%s-%s", functionName, pluginName)
					}
				}
			}
		}
	}

	return annotations
}
