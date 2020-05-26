package get

import (
	"fmt"
)

const helmVersion = "v3.1.2"

func getHelmUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz", helmVersion, os, arch)
}

func DownloadHelm() error {
	helmUrl := getHelmUrl()

	if err := download(helmUrl, "helm", false); err != nil {
		return err
	}

	return nil
}
