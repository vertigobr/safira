// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	s "github.com/vertigobr/safira/pkg/stack"
	"strings"
)

type functionSpec struct {
	Name     string            `yaml:"name,omitempty"`
	Image    string            `yaml:"image,omitempty"`
	Labels   map[string]string `yaml:"labels,omitempty"`
	Limits   cpuMemory         `yaml:"limits,omitempty"`
	Requests cpuMemory         `yaml:"requests,omitempty"`
}

type cpuMemory struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

func (k *K8sYaml) MountFunction(functionName, namespace string) error {
	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	scaleMin, scaleMax := getScaleConfig(stack, functionName)
	*k = K8sYaml{
		ApiVersion: "openfaas.com/v1",
		Kind:       "Function",
		Metadata: metadata{
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
