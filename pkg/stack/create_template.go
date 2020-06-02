// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func CreateTemplate(functionName, templateName, handlerFile, image string) error {
	stack := Stack{
		Version:   "1.0",
		Provider:  Provider{
			Name: "openfaas",
			GatewayURL: "http://gateway.ipaas.localdomain:8080",
		},
		Functions: map[string]Function{
			functionName: {
				Name: functionName,
				Template: templateName,
				Handler: handlerFile,
				Image: image,
			},
		},
	}

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return err
	}

	if err := utils.CreateYamlFile("stack.yml", yamlBytes, true); err != nil {
		return err
	}

	return nil
}
