// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package kong

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

type kongIngressPlugin struct {
	ApiVersion string            `yaml:"apiVersion,omitempty"`
	Kind       string            `yaml:"kind,omitempty"`
	Metadata   map[string]string `yaml:"metadata,omitempty"`
	Upstream   map[string]string `yaml:"upstream,omitempty"`
	Proxy      map[string]string `yaml:"proxy,omitempty"`
	Route      map[string]string `yaml:"route,omitempty"`
}


func createIngress(pluginName, kongpluginFolder string) error {
	kp := kongIngressPlugin{
		ApiVersion: "configuration.konghq.com/v1",
		Kind: "KongIngress",
		Metadata: map[string]string{
			"name": pluginName,
		},
		Upstream: map[string]string{
			"<key>": "<value>",
		},
		Proxy: map[string]string{
			"<key>": "<value>",
		},
		Route: map[string]string{
			"<key>": "<value>",
		},
	}

	folder := fmt.Sprintf("%s/%s.yml", kongpluginFolder, pluginName)

	yamlBytes, err := y.Marshal(&kp)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", folder, err.Error())
	}

	if err := utils.CreateYamlFile(folder, yamlBytes, true); err != nil {
		return err
	}

	return nil
}
