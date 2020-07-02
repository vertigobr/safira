// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package kong

import (
	"fmt"
	"os"

	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func Create(assetName, assetType string) error {
	kongpluginFolder := GetKongFolderName()

	if _, err := os.Stat(kongpluginFolder); err != nil {
		if err = os.MkdirAll(kongpluginFolder, 0700); err != nil {
			return err
		}

		stack, err := s.LoadStackFile()
		if err != nil {
			return err
		}

		stack.KongAssetsEnabled = true
		yamlBytes, err := y.Marshal(&stack)
		if err != nil {
			return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", s.GetYamlFileName(), err.Error())
		}

		if err := utils.CreateYamlFile(s.GetYamlFileName(), yamlBytes, true); err != nil {
			return err
		}
	}

	if _, err := os.Stat(kongpluginFolder + "/" + assetName + ".yml"); err == nil {
		return fmt.Errorf("\nNome de plugin já utilizado!")
	}

	switch assetType {
	case "plugin":
		if err := createPlugin(assetName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "cluster-plugin":
		if err := createClusterPlugin(assetName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "ingress":
		if err := createIngress(assetName, kongpluginFolder); err != nil {
			return err
		}
		break
	case "consumer":
		if err := createConsumer(assetName, kongpluginFolder); err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf(fmt.Sprintf("Tipo %s inválido!", assetType))
	}

	return nil
}

func GetKongFolderName() string {
	return "kong"
}
