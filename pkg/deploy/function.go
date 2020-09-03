// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"strings"

	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/git"
	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
)

type functionSpec struct {
	Name         string                 `yaml:"name,omitempty"`
	Image        string                 `yaml:"image,omitempty"`
	Labels       map[string]string      `yaml:"labels,omitempty"`
	Environments map[string]interface{} `yaml:"environment,omitempty"`
	Limits       cpuMemory              `yaml:"limits,omitempty"`
	Requests     cpuMemory              `yaml:"requests,omitempty"`
}

type cpuMemory struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

func (k *K8sYaml) MountFunction(functionName, namespace, env string, useSha bool) error {
	stack, err := s.LoadStackFile(env)
	if err != nil {
		return err
	}

	scaleMin, scaleMax := getScaleConfig(stack, functionName)
	cpuLimits, memoryLimits := getLimitsConfig(stack, functionName)
	cpuRequests, memoryRequests := getRequestsConfig(stack, functionName)
	environments := getFunctionEnvironment(stack, functionName)

	image := stack.Functions[functionName].Image
	if useSha {
		imageWithCommitSha, _ := git.GetImageWithCommitSha(image)
		if len(imageWithCommitSha) > 0 {
			image = imageWithCommitSha
		}
	}

	repoName, err := utils.GetCurrentFolder()
	if err != nil {
		return err
	}

	*k = K8sYaml{
		ApiVersion: "openfaas.com/v1",
		Kind:       "Function",
		Metadata: metadata{
			Name:      functionName,
			Namespace: namespace,
			Annotations: map[string]string{
				"safira.io/repository": repoName,
				"safira.io/function":   functionName,
			},
		},
		Spec: functionSpec{
			Name:  functionName,
			Image: image,
			Labels: map[string]string{
				"com.openfaas.scale.min": scaleMin,
				"com.openfaas.scale.max": scaleMax,
				"function":               functionName,
			},
			Environments: environments,
			Limits: cpuMemory{
				Cpu:    cpuLimits,
				Memory: memoryLimits,
			},
			Requests: cpuMemory{
				Cpu:    cpuRequests,
				Memory: memoryRequests,
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
		Command: config.GetKubectlPath(),
		Args: []string{
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

	min = compareFuncStackValue(minFunction, minStack, "0")
	max = compareFuncStackValue(maxFunction, maxStack, "5")

	return
}

func getLimitsConfig(stack *s.Stack, functionName string) (cpu, memory string) {
	cpuFunction := stack.Functions[functionName].FunctionConfig.Limits.Cpu
	memoryFunction := stack.Functions[functionName].FunctionConfig.Limits.Memory
	cpuStack := stack.StackConfig.Limits.Cpu
	memoryStack := stack.StackConfig.Limits.Memory

	cpu = compareFuncStackValue(cpuFunction, cpuStack, "")
	memory = compareFuncStackValue(memoryFunction, memoryStack, "")

	return
}

func getRequestsConfig(stack *s.Stack, functionName string) (cpu, memory string) {
	cpuFunction := stack.Functions[functionName].FunctionConfig.Requests.Cpu
	memoryFunction := stack.Functions[functionName].FunctionConfig.Requests.Memory
	cpuStack := stack.StackConfig.Requests.Cpu
	memoryStack := stack.StackConfig.Requests.Memory

	cpu = compareFuncStackValue(cpuFunction, cpuStack, "")
	memory = compareFuncStackValue(memoryFunction, memoryStack, "")

	return
}

func getFunctionEnvironment(stack *s.Stack, functionName string) map[string]interface{} {
	envFunction := stack.Functions[functionName].FunctionConfig.Environments
	envStack := stack.StackConfig.Environments

	if len(envFunction) > 0 {
		return envFunction
	} else if len(envStack) > 0 {
		return envStack
	}

	var m map[string]interface{}
	return m
}

func compareFuncStackValue(functionValue, stackValue, defaultValue string) (value string) {
	if len(functionValue) > 0 {
		value = functionValue
	} else if len(stackValue) > 0 {
		value = stackValue
	} else {
		value = defaultValue
	}

	return
}
