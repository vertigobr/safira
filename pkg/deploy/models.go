// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"

	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	"gopkg.in/gookit/color.v1"
	y "gopkg.in/yaml.v2"
)

// Yaml file structure
type K8sYaml struct {
	ApiVersion string            `yaml:"apiVersion,omitempty"`
	Kind       string            `yaml:"kind,omitempty"`
	Metadata   metadata          `yaml:"metadata,omitempty"`
	Spec       interface{}       `yaml:"spec,omitempty"`
	Data       map[string]string `yaml:"data,omitempty"`

	// Kong Plugin
	Config     map[string]interface{} `yaml:"config,omitempty"`
	ConfigFrom s.ConfigFromPlugin     `yaml:"configFrom,omitempty"`
	Plugin     string                 `yaml:"plugin,omitempty"`
}

type metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

func (k *K8sYaml) CreateYamlFile(fileName string) error {
	yamlBytes, err := y.Marshal(&k)
	if err != nil {
		return fmt.Errorf("%s Error reading the file: %s", color.Red.Text("[!]"), fileName)
	}

	if err := utils.CreateYamlFile(fileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}
