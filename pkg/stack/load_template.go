// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"fmt"
	"io/ioutil"

	y "gopkg.in/yaml.v2"
)

func LoadStackFile(file string) (*Stack, error) {
	var stack Stack
	yamlBytes, err := ParseYAMLForLanguageTemplate(file)
	if err != nil {
		return nil, err
	}

	err = y.Unmarshal(yamlBytes, &stack)
	if err != nil {
		fmt.Printf("Error with YAML file\n")
		return nil, err
	}

	return &stack, nil
}

func ParseYAMLForLanguageTemplate(file string) (fileData []byte, err error) {
	fileData, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return
}