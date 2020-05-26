package get

import (
	"fmt"
)

const faasVersion = "0.12.2"

func getFaasCliUrl() string {
	var suffix string
	arch := getARCH()
	os := getOS()

	if os == "darwin" {
		suffix = "-darwin"
	} else if os == "mingw" || os == "windows" {
		suffix = ".exe"
	} else {
		if arch == "aarch64" {
			suffix = "-arm64"
		} else if arch == "armv6l" || arch == "armv7l" {
			suffix = "-armhf"
		}
	}

	return fmt.Sprintf("https://github.com/openfaas/faas-cli/releases/download/%s/faas-cli%s", faasVersion, suffix)
}

func DownloadFaasCli() error {
	faasCliUrl := getFaasCliUrl()

	if err := download(faasCliUrl, "faas-cli", true); err != nil {
		return err
	}

	return nil
}
