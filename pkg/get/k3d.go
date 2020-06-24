// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import "fmt"

const k3dVersion = "v1.7.0"

func getK3dUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://github.com/rancher/k3d/releases/download/%s/k3d-%s-%s", k3dVersion, os, arch)
}

func DownloadK3d() error {
	k3dUrl := getK3dUrl()

	if err := downloadBinary(k3dUrl, "k3d", true); err != nil {
		return err
	}

	return nil
}
