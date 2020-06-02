// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"
	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

type ingress struct {
	ApiVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   ingressMetadata `yaml:"metadata"`
	Spec       ingressSpec     `yaml:"spec"`
}

type ingressMetadata struct {
	Name string `yaml:"name"`
}

type ingressSpec struct {
	Rules []rule `yaml:"rules"`
}

type rule struct {
	Host string `yaml:"host"`
	Http http   `yaml:"http"`
}

type http struct {
	Paths []path `yaml:"paths"`
}

type path struct {
	Path    string  `yaml:"path"`
	Backend backend `yaml:"backend"`
}

type backend struct {
	ServiceName string `yaml:"serviceName"`
	ServicePort int    `yaml:"servicePort"`
}

func CreateYamlIngress(fileName, functionName string) error {
	stack, err := s.LoadStackFile("./stack.yml")
	if err != nil {
		return err
	}

	gateway, port, err := getGatewayPort(stack.Provider.GatewayURL)
	if err != nil {
		return err
	}

	ingress := ingress{
		ApiVersion: "extensions/v1beta1",
		Kind:       "Ingress",
		Metadata: ingressMetadata{
			Name: functionName,
		},
		Spec: ingressSpec{
			Rules: []rule{
				{
					Host: gateway,
					Http: http{
						Paths: []path{
							{
								Path: "/function/" + functionName,
								Backend: backend{
									ServiceName: functionName,
									ServicePort: port,
								},
							},
						},
					},
				},
			},
		},
	}

	yamlBytes, err := y.Marshal(&ingress)
	if err != nil {
		return err
	}

	if err := utils.CreateYamlFile(fileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}

func getGatewayPort(url string) (gateway string, port int, err error) {
	if strings.Index(url, "http://") != -1 {
		gateway = strings.Trim(url, "http://")
	} else if strings.Index(url, "https://") != -1 {
		gateway = strings.Trim(url, "https://")
	}

	s := strings.Split(gateway, ":")
	gateway = s[0]
	port, err = strconv.Atoi(s[1])
	if err != nil {
		return "", 0, fmt.Errorf("stack: url do gateway inválida")
	}

	return
}
