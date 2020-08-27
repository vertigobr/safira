// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package config

import (
	"fmt"
	"os"
	"os/user"
	p "path"
	"strconv"

	"gopkg.in/gookit/color.v1"
)

func GetSafiraDir() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/.safira/", home)
}

func initUserDir(folder string) (string, error) {
	safiraDir := GetSafiraDir()

	path := p.Join(safiraDir, folder)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("%s Error to create a folder: %s", color.Red.Text("[!]"), folder)
	}

	return path, nil
}

func CreateInBinDir() (string, error) {
	path, err := initUserDir("/bin/")
	if err != nil {
		return "", err
	}

	sudoUser := os.Getenv("SUDO_USER")
	if len(sudoUser) > 0 {
		u, _ := user.Lookup(sudoUser)
		Uid, _ := strconv.Atoi(u.Uid)
		Gid, _ := strconv.Atoi(u.Gid)

		if err := os.Chown(path, Uid, Gid); err != nil {
			return "", fmt.Errorf("%s Error when changing the owner of the root folder for the user: %s", color.Red.Text("[!]"), err.Error())
		}

		safiraFolder := GetSafiraDir()
		if err := os.Chown(safiraFolder, Uid, Gid); err != nil {
			return "", fmt.Errorf("%s Error when changing the owner of the root folder for the user: %s", color.Red.Text("[!]"), err.Error())
		}
	}

	return path, err
}

func SetKubeconfig(clusterName string) error {
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME")+"/.config/k3d/"+clusterName+"/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("%s Unable to export KUBECONFIG", color.Red.Text("[!]"))
	}

	return nil
}

func GetKubeconfig() string {
	return os.Getenv("KUBECONFIG")
}
