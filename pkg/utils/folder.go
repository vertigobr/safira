package utils

import (
	"fmt"
	"os"
	"path"
)

func GetCurrentFolder() (string, error) {
	pt, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current folder")
	}

	_, folder := path.Split(pt)

	return folder, nil
}
