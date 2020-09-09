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

	// Path
	Path string `yaml:"path,omitempty"`

	// Config
	FunctionConfig Config `yaml:"config,omitempty"`

	// Kong Plugins
	Plugins map[string]Plugin `yaml:"plugins,omitempty"`
}

// Config apply one or all functions in stack.yaml
type Config struct {
	Build  Build  `yaml:"build,omitempty"`
	Deploy Deploy `yaml:"deploy,omitempty"`
	Scale  struct {
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"scale,omitempty"`

	// Resource limit that the function will use
	Limits CpuMemory `yaml:"limits,omitempty"`

	// Minimum resource required that the function will use
	Requests CpuMemory `yaml:"requests,omitempty"`

	// Environment variables
	Environments map[string]interface{} `yaml:"environments,omitempty"`
}

type Build struct {
	Enabled *bool             `yaml:"enabled,omitempty"`
	UseSha  bool              `yaml:"useSha,omitempty"`
	Args    map[string]string `yaml:"args,omitempty"`
}

type Deploy struct {
	//Enabled *bool  `yaml:"enabled,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
	Suffix string `yaml:"suffix,omitempty"`
}

type CpuMemory struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type Plugin struct {
	// Plugin name
	Name string `yaml:"-"`

	// Plugin type
	Type string `yaml:"type,omitempty"`

	Global string `yaml:"global,omitempty"`

	Config map[string]interface{} `yaml:"config,omitempty"`

	ConfigFrom ConfigFromPlugin `yaml:"configFrom,omitempty"`
}

type ConfigFromPlugin struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef,omitempty"`
}

type SecretKeyRef struct {
	Name string `yaml:"name,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

// Stack root level YAML file to define FaaS function-set
type Stack struct {
	Version     string              `yaml:"version,omitempty"`
	Provider    Provider            `yaml:"provider,omitempty"`
	Hostname    string              `yaml:"hostname,omitempty"`
	Swagger     Swagger             `yaml:"swagger,omitempty"`
	Functions   map[string]Function `yaml:"functions,omitempty"`
	StackConfig Config              `yaml:"config,omitempty"`
	Custom      []string            `yaml:"custom,omitempty"`
}

type Swagger struct {
	File string `yaml:"file,omitempty"`
}

func GetStackFileName() string {
	return "stack.yml"
}
