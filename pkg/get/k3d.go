package get

import (
	"fmt"
)

const k3dVersion = "v1.7.0"

func getK3dUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://github.com/rancher/k3d/releases/download/%s/k3d-%s-%s", k3dVersion, os, arch)
}

func DownloadK3d() error {
	k3dUrl := getK3dUrl()

	if err := download(k3dUrl, "k3d", true); err != nil {
		return err
	}

	return nil
}
