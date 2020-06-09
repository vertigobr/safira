// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
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

func CreateYamlFunction(fileName, functionName, namespace string) error {
	stack, err := s.LoadStackFile()
	scaleMin, scaleMax := getScaleConfig(stack, functionName)
	
	function := function{
		ApiVersion: "openfaas.com/v1",
		Kind:       "Function",
		Metadata: functionMetadata{
			Name:      functionName,
			Namespace: namespace,
		},
		Spec: functionSpec{
			Name:  functionName,
			Image: stack.Functions[functionName].Image,
			Labels: map[string]string{
				"com.openfaas.scale.min": scaleMin,
				"com.openfaas.scale.max": scaleMax,
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

func CheckFunction(clusterName, functionName, namespace string) (bool, error) {
	err := config.SetKubeconfig(clusterName)
	if err != nil {
		return false, err
	}

	taskCheckFunction := execute.Task{
		Command:     config.GetKubectlPath(),
		Args:        []string{
			"get", "deployments", "-n", namespace,
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

func getScaleConfig(stack *s.Stack, functionName string) (min, max string) {
	minFunction := stack.Functions[functionName].FunctionConfig.Scale.Min
	maxFunction := stack.Functions[functionName].FunctionConfig.Scale.Max
	minStack := stack.StackConfig.Scale.Min
	maxStack := stack.StackConfig.Scale.Max

	if len(minFunction) > 0 {
		min = minFunction
	} else if len(minStack) > 0 {
		min = minStack
	} else {
		min = "0"
	}

	if len(maxFunction) > 0 {
		max = maxFunction
	} else if len(maxStack) > 0 {
		max = maxStack
	} else {
		max = "5"
	}

	return
}
