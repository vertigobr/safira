package config

import (
	"fmt"
	"os"
	"os/exec"
	p "path"
)

func GetUserDir() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/.safira/", home)
}

func initUserDir(root string) (string, error) {
	safiraDir := GetUserDir()

	if len(safiraDir) == 0{
		return safiraDir, fmt.Errorf("variável HOME não encontrada")
	}

	path := p.Join(safiraDir, root)
	if err := os.MkdirAll(path, 0700); err != nil {
		return path, err
	}

	return path, nil
}

func CreateInBinDir() (string, error) {
	return initUserDir("/bin/")
}

func CreateInTemplateDir() (string, error) {
	return initUserDir("/template/")
}

func ExistsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", GetUserDir(), binary))
	exists = path != ""
	return
}
