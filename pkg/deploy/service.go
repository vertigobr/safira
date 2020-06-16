// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	s "github.com/vertigobr/safira/pkg/stack"
)

type serviceSpec struct {
	Type         string        `yaml:"type"`
	ExternalName string        `yaml:"externalName"`
	Ports        []servicePort `yaml:"ports"`
}

type servicePort struct {
	Port int `yaml:"port"`
}

func (k *K8sYaml) MountService(serviceName, hostnameFlag string) error {
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
		Spec: serviceSpec{
			Type: "ExternalName",
			ExternalName: "gateway",
			Ports: []servicePort{
				{
					Port: p,
				},
			},
		},
	}

	return nil
}
