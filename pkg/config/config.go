// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package config

import (
	"fmt"
	"os"
	"os/user"
	p "path"
	"strconv"
)

func GetUserDir() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/.safira/", home)
}

func initUserDir(folder string) (string, error) {
	safiraDir := GetUserDir()

	if len(safiraDir) <= 16 {
		return "", fmt.Errorf("variável SUDO_USER não encontrada")
	}

	path := p.Join(safiraDir, folder)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("error ao criar pasta: %s", folder)
	}

	return path, nil
}

func CreateInBinDir() (string, error) {
	path, err := initUserDir("/bin/")
	if err != nil {
		return "", err
	}

	sudoUser := os.Getenv("SUDO_USER")
	u, _ := user.Lookup(sudoUser)
	Uid, _ := strconv.Atoi(u.Uid)
	Gid, _ := strconv.Atoi(u.Gid)

	if err := os.Chown(path, Uid, Gid); err != nil {
		return "", fmt.Errorf("error ao mudar o dono da pasta de root para o usuário: %s", err.Error())
	}

	safiraFolder := GetUserDir()
	if err := os.Chown(safiraFolder, Uid, Gid); err != nil {
		return "", fmt.Errorf("error ao mudar o dono da pasta de root para o usuário: %s", err.Error())
	}

	return path, err
}

func SetKubeconfig(clusterName string) error {
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível criar a variável de ambiente: KUBECONFIG")
	}

	return nil
}

func GetKubeconfig() string {
	return os.Getenv("KUBECONFIG")
}
