// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
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

func CreateYamlFunction(fileName string) error {
	if err := readFileEnv(); err != nil {
		return err
	}

	projectName, imageName, err := getFunctionEnvs()
	if err != nil {
		return nil
	}
	
	function := function{
		ApiVersion: "openfaas.com/v1",
		Kind:       "Function",
		Metadata: functionMetadata{
			Name:      projectName,
			Namespace: GetNamespaceFunction(),
		},
		Spec: functionSpec{
			Name:  projectName,
			Image: imageName,
			Labels: map[string]string{
				"com.openfaas.scale.min": "3",
				"com.openfaas.scale.max": "5",
				"function": projectName,
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
		return err
	}

	if err := createYamlFile(fileName, yamlBytes); err != nil {
		return err
	}

	return nil
}

func getFunctionEnvs() (string, string, error) {
	projectName, err := GetProjectName()
	if err != nil {
		return "", "", err
	}

	imageName, err := getImageName()
	if err != nil {
		return "", "", err
	}

	return projectName, imageName, nil
}

func GetNamespaceFunction() string {
	return "openfaas-fn"
}

func CheckFunction(clusterName string) (bool, error) {
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

	projectName, err := GetProjectName()

	if strings.Contains(res.Stdout, projectName) {
		return true, nil
	}

	return false, nil
}
