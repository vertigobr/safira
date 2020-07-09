// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

// Provider for the FaaS set of functions
type Provider struct {
	Name       string `yaml:"name"`
	GatewayURL string `yaml:"gateway"`
}

// Function as deployed or built on FaaS
type Function struct {
	// Name of deployed function
	Name string `yaml:"-"`

	// Lang/Template name
	Lang string `yaml:"lang"`

	// Handler Local folder to use for function
	Handler string `yaml:"handler"`

	// Docker image name
	Image string `yaml:"image"`

	FunctionConfig Config `yaml:"config,omitempty"`
}

// Config apply one or all functions in stack.yaml
type Config struct {
	BuildArgs    map[string]string      `yaml:"buildArgs,omitempty"`
	Scale        struct{
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"scale,omitempty"`

	// Resource limit that the function will use
	Limits       CpuMemory              `yaml:"limits,omitempty"`

	// Minimum resource required that the function will use
	Requests     CpuMemory              `yaml:"requests,omitempty"`

	// Environment variables
	Environments map[string]interface{} `yaml:"environments,omitempty"`
}

type CpuMemory struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// Stack root level YAML file to define FaaS function-set
type Stack struct {
	Version            string              `yaml:"version,omitempty"`
	Provider           Provider            `yaml:"provider,omitempty"`
	Hostname           string              `yaml:"hostname,omitempty"`
	Functions          map[string]Function `yaml:"functions,omitempty"`
	StackConfig        Config              `yaml:"config,omitempty"`
	Custom             []string            `yaml:"custom,omitempty"`
	KongAssetsEnabled  bool                `yaml:"kongAssetsEnabled,omitempty"`
}

func GetYamlFileName() string {
	return "stack.yml"
}
