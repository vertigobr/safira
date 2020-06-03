// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"
	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

type service struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   serviceMetadata `yaml:"metadata"`
	Spec       serviceSpec     `yaml:"spec"`
}

type serviceMetadata struct {
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}

type serviceSpec struct {
	Type         string            `yaml:"type"`
	ExternalName string            `yaml:"externalName"`
	Ports        []port            `yaml:"ports"`
}

type port struct {
	Port int `yaml:"port"`
}

func CreateYamlService(fileName, functionName string) error {
	stack, err := s.LoadStackFile("./stack.yml")
	if err != nil {
		return err
	}

	_, p, err := getGatewayPort(stack.Provider.GatewayURL)
	if err != nil {
		return err
	}

	service := service{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: serviceMetadata{
			Name:   functionName,
			Labels: map[string]string{
				"app": functionName,
			},
			Annotations: map[string]string{
				"konghq.com/plugins": "prometheus",
			},
		},
		Spec: serviceSpec{
			Type: "ExternalName",
			ExternalName: "gateway",
			Ports: []port{
				{
					Port: p,
				},
			},
		},
	}

	yamlBytes, err := y.Marshal(&service)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", fileName, err.Error())
	}

	if err := utils.CreateYamlFile(fileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}
