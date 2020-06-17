// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
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

func (k *K8sYaml) MountService(serviceName, hostnameFlag string, function bool) error {
	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	var p int
	if len(hostnameFlag) > 1 {
		_, p, err = getGatewayPort(hostnameFlag)
	} else {
		_, p, err = getGatewayPort(stack.Hostname)
	}

	if err != nil {
		return err
	}

	var spec serviceSpec
	if function {
		spec = serviceSpec{
			Type: "ExternalName",
			ExternalName: "gateway",
			Ports: []servicePort{
				{
					Port: p,
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

	*k = K8sYaml{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: metadata{
			Name:   serviceName,
			Labels: map[string]string{
				"app": serviceName,
			},
			Annotations: map[string]string{
				"konghq.com/plugins": "prometheus",
			},
		},
		Spec: spec,
	}

	return nil
}
