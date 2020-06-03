// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"fmt"
	"io/ioutil"

	y "gopkg.in/yaml.v2"
)

func LoadStackFile(fileName string) (*Stack, error) {
	var stack Stack
	yamlBytes, err := ParseYAMLForLanguageTemplate(fileName)
	if err != nil {
		return nil, err
	}

	err = y.Unmarshal(yamlBytes, &stack)
	if err != nil {
		return nil, fmt.Errorf("error ao executar o unmarshalling para o arquivo %s: %s", fileName, err.Error())
	}

	return &stack, nil
}

func ParseYAMLForLanguageTemplate(fileName string) (fileData []byte, err error) {
	fileData, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error ao let o arquivo %s: %s", fileName, err.Error())
	}

	return
}