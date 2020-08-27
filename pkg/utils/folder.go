// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package utils

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/gookit/color.v1"
)

func GetCurrentFolder() (string, error) {
	pt, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%s Error getting path to current folder", color.Red.Text("[!]"))
	}

	_, folder := path.Split(pt)

	return folder, nil
}
