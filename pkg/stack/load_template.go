// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/gookit/color.v1"
	y "gopkg.in/yaml.v2"
)

func LoadStackFile(envFile string) (*Stack, error) {
	var stack Stack
	yamlBytes, err := ParseYAMLForLanguageTemplate(GetStackFileName())
	if err != nil {
		return nil, err
	}

	err = y.Unmarshal(yamlBytes, &stack)
	if err != nil {
		return nil, fmt.Errorf("%s Error when unmarshalling the %s file", color.Red.Text("[!]"), GetStackFileName())
	}

	if len(envFile) > 0 {
		var envStack Stack
		yamlBytes, err = ParseYAMLForLanguageTemplate(envFile)
		if err != nil {
			return nil, err
		}

		err = y.Unmarshal(yamlBytes, &envStack)
		if err != nil {
			return nil, fmt.Errorf("%s Error when unmarshalling the %s file", color.Red.Text("[!]"), GetStackFileName())
		}

		if err := prepareStack(&stack, &envStack); err != nil {
			return nil, err
		}
	}

	return &stack, nil
}

func ParseYAMLForLanguageTemplate(fileName string) (fileData []byte, err error) {
	fileData, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s Error reading %s file", color.Red.Text("[!]"), fileName)
	}

	return
}
