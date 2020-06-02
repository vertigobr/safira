// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import "fmt"

const kubectlVersion = "v1.18.0"

func getKubectlUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", kubectlVersion, os, arch)
}

func DownloadKubectl() error {
	kubectlUrl := getKubectlUrl()

	if err := download(kubectlUrl, "kubectl", true); err != nil {
		return err
	}

	return nil
}
