// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
)

func readFileEnv() error {
	if err := gotenv.Load(); err != nil {
		return err
	}

	return nil
}

func GetProjectName() (string, error) {
	s := os.Getenv("PROJECT_NAME")
	if len(s) == 0 {
		return "", fmt.Errorf("variável PROJECT_NAME não encontrada no arquivo .env")
	}

	return s, nil
}

func getImageName() (string, error) {
	s := os.Getenv("IMAGE")
	if len(s) == 0 {
		return "", fmt.Errorf("variável IMAGE não encontrada no arquivo .env")
	}

	return s, nil
}

func getDomain() (string, error) {
	s := os.Getenv("DOMAIN")
	if len(s) == 0 {
		return "", fmt.Errorf("variável DOMAIN não encontrada no arquivo .env")
	}

	return s, nil
}

func getPort() (int, error) {
	s := os.Getenv("PORT")
	if len(s) == 0 {
		return 0, fmt.Errorf("variável PORT não encontrada no arquivo .env")
	}

	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return port, nil
}
