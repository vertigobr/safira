// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package utils

import (
	"os"
	"strings"
)

func CreateYamlFile(fileName string, bytes []byte, clearFile bool) error {
	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		fileName = fileName + ".yml"
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if clearFile {
		f.Truncate(0)
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}