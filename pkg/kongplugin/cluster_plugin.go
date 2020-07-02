// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package kongplugin

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func createYamlClusterPlugin(pluginName, kongpluginFolder string) error {
	kp := kongPlugin{
		ApiVersion: "configuration.konghq.com/v1",
		Kind: "KongClusterPlugin",
		Metadata: map[string]string{
			"name": pluginName,
		},
		Config: map[string]string{
			"<key>": "<value>",
		},
		Plugin: "<plugin name>",
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
