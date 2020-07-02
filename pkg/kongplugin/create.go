// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package kongplugin

import (
	"fmt"
	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
	"os"
)

func CreatePlugin(pluginName, typeFlag string) error {
	kongpluginFolder := GetKongPluginFolderName()

	if _, err := os.Stat(kongpluginFolder); err != nil {
		if err = os.MkdirAll(kongpluginFolder, 0700); err != nil {
			return err
		}

		stack, err := s.LoadStackFile()
		if err != nil {
			return err
		}

		stack.KongPluginsEnabled = true
		yamlBytes, err := y.Marshal(&stack)
		if err != nil {
			return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", s.GetYamlFileName(), err.Error())
		}

		if err := utils.CreateYamlFile(s.GetYamlFileName(), yamlBytes, true); err != nil {
			return err
		}
	}

	switch typeFlag {
	case "plugin":
		if err := createYamlPlugin(pluginName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "clusterPlugin":
		if err := createYamlClusterPlugin(pluginName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "ingress":
		if err := createYamlIngressPlugin(pluginName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "consumer":
		if err := createYamlConsumer(pluginName, kongpluginFolder); err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf("valor da flag type é inválido")
	}

	return nil
}

func GetKongPluginFolderName() string {
	return "kongplugin"
}
