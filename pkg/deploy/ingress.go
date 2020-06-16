// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"
	s "github.com/vertigobr/safira/pkg/stack"
	"strconv"
	"strings"
)

type ingressSpec struct {
	Rules []ingressRule `yaml:"rules,omitempty"`
}

type ingressRule struct {
	Host string      `yaml:"host,omitempty"`
	Http ingressHttp `yaml:"http,omitempty"`
}

type ingressHttp struct {
	Paths []ingressPath `yaml:"paths,omitempty"`
}

type ingressPath struct {
	Path    string         `yaml:"path,omitempty"`
	Backend ingressBackend `yaml:"backend,omitempty"`
}

type ingressBackend struct {
	ServiceName string `yaml:"serviceName,omitempty"`
	ServicePort int    `yaml:"servicePort,omitempty"`
}

func (k *K8sYaml) MountIngress(ingressName, serviceName, path, hostnameFlag string) error {
	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	var port int
	var gateway string
	if len(hostnameFlag) > 1 {
		gateway, port, err = getGatewayPort(hostnameFlag)
	} else {
		gateway, port, err = getGatewayPort(stack.Hostname)
	}

	if err != nil {
		return err
	}

	*k = K8sYaml{
		ApiVersion: "extensions/v1beta1",
		Kind:       "Ingress",
		Metadata: metadata{
			Name: ingressName,
		},
		Spec: ingressSpec{
			Rules: []ingressRule{
				{
					Host: gateway,
					Http: ingressHttp{
						Paths: []ingressPath{
							{
								Path: "/function/" + path,
								Backend: ingressBackend{
									ServiceName: serviceName,
									ServicePort: port,
								},
							},
						},
					},
				},
			},
		},
	}

	return nil
}

func getGatewayPort(url string) (gateway string, port int, err error) {
	if strings.Index(url, "http://") != -1 {
		gateway = strings.Trim(url, "http://")
	} else if strings.Index(url, "https://") != -1 {
		gateway = strings.Trim(url, "https://")
	} else {
		gateway = url
	}

	s := strings.Split(gateway, ":")
	gateway = s[0]
	port, err = strconv.Atoi(s[1])
	if err != nil {
		return "", 0, fmt.Errorf("error ao pegar a porta do hostname: %s", err.Error())
	}

	return
}
