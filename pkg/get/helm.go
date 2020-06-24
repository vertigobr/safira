// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import "fmt"

const helmVersion = "v3.1.2"

func getHelmUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz", helmVersion, os, arch)
}

func DownloadHelm() error {
	helmUrl := getHelmUrl()

	if err := downloadBinary(helmUrl, "helm", false); err != nil {
		return err
	}

	return nil
}
