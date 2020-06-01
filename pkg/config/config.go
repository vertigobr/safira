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
		return "", err
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
		return "", err
	}

	safiraFolder := GetUserDir()
	if err := os.Chown(safiraFolder, Uid, Gid); err != nil {
		return "", err
	}

	return path, err
}

//func CreateInTemplateDir() (string, error) {
//	return initUserDir("/template/")
//}

func SetKubeconfig(clusterName string) error {
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível criar a variável de ambiente: KUBECONFIG")
	}

	return nil
}
