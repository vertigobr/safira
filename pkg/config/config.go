package config

import (
	"fmt"
	"os"
	p "path"
)

func getUserDir() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/.safira/", home)
}

func initUserDir(root, folder string) (string, error) {
	safiraDir := getUserDir()

	if len(safiraDir) == 0{
		return safiraDir, fmt.Errorf("variável HOME não encontrada")
	}

	path := p.Join(safiraDir, root)
	if err := os.MkdirAll(path, 0700); err != nil {
		return path, err
	}

	subPath := p.Join(path, folder)
	if err := os.MkdirAll(subPath, 0700); err != nil {
		return subPath, err
	}

	return subPath, nil
}

func CreateInBinDir(folder string) (string, error) {
	return initUserDir("/bin/", folder)
}

func CreateInTemplateDir(folder string) (string, error) {
	return initUserDir("/template/", folder)
}
