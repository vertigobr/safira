// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
	"strings"
)

type function struct {
	ApiVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   functionMetadata `yaml:"metadata"`
	Spec       functionSpec     `yaml:"spec"`
}

type functionMetadata struct {
	Name        string `yaml:"name"`
	Namespace   string `yaml:"namespace"`
}

type functionSpec struct {
	Name         string            `yaml:"name"`
	Image        string            `yaml:"image"`
	Labels       map[string]string `yaml:"labels"`
	Limits       cpuMemory         `yaml:"limits"`
	Requests     cpuMemory         `yaml:"requests"`
}

type cpuMemory struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

func CreateYamlFunction(fileName, functionName string) error {
	stack, err := s.LoadStackFile("./stack.yml")
	
	function := function{
		ApiVersion: "openfaas.com/v1",
		Kind:       "Function",
		Metadata: functionMetadata{
			Name:      functionName,
			Namespace: GetNamespaceFunction(),
		},
		Spec: functionSpec{
			Name:  functionName,
			Image: stack.Functions[functionName].Image,
			Labels: map[string]string{
				"com.openfaas.scale.min": "3",
				"com.openfaas.scale.max": "5",
				"function": functionName,
			},
			Limits: cpuMemory{
				Cpu:    "200m",
				Memory: "256Mi",
			},
			Requests: cpuMemory{
				Cpu:    "10m",
				Memory: "128Mi",
			},
		},
	}

	yamlBytes, err := y.Marshal(&function)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", fileName, err.Error())
	}

	if err := utils.CreateYamlFile(fileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}

func GetNamespaceFunction() string {
	return "openfaas-fn"
}

func CheckFunction(clusterName, functionName string) (bool, error) {
	err := config.SetKubeconfig(clusterName)
	if err != nil {
		return false, err
	}

	taskCheckFunction := execute.Task{
		Command:     config.GetKubectlPath(),
		Args:        []string{
			"get", "deployments", "-n", GetNamespaceFunction(),
		},
		StreamStdio:  false,
		PrintCommand: false,
	}

	res, err := taskCheckFunction.Execute()
	if err != nil {
		return false, err
	}

	if res.ExitCode != 0 {
		return false, nil
	}

	if strings.Contains(res.Stdout, functionName) {
		return true, nil
	}

	return false, nil
}
