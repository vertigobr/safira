// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func CreateTemplate(function Function) error {
	stack := Stack{
		Version:   "1.0",
		Provider:  Provider{
			Name:       "openfaas",
			GatewayURL: "http://gateway.ipaas.localdomain:8080",
		},
		Hostname: "ipaas.localdomain:8080",
		Functions: map[string]Function{
			function.Name: {
				Name:    function.Name,
				Lang:    function.Lang,
				Handler: function.Handler,
				Image:   function.Image,
			},
		},
	}

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo stack.yml: %s", err.Error())
	}

	if err := utils.CreateYamlFile(GetYamlFileName(), yamlBytes, true); err != nil {
		return err
	}

	return nil
}
