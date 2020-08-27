// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/gookit/color.v1"
)

func CreateYamlFile(fileName string, bytes []byte, clearFile bool) error {
	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		fileName = fileName + ".yml"
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("%s Error opening the %s file", color.Red.Text("[!]"), fileName)
	}
	defer f.Close()

	if clearFile {
		f.Truncate(0)
	}

	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("%s Error writing the %s file", color.Red.Text("[!]"), fileName)
	}

	return nil
}

func AppendYamlFile(fileName string, bytes []byte) error {
	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		fileName = fileName + ".yml"
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("%s Error opening the %s file", color.Red.Text("[!]"), fileName)
	}
	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("%s Error writing the %s file", color.Red.Text("[!]"), fileName)
	}

	return nil
}
