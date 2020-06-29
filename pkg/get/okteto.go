// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
)

const oktetoVersion = "1.8.12"

func getOktetoUrl() string {
	var suffix string
	arch := getARCH()
	os := getOS()

	if os == "darwin" {
		suffix = "-Darwin-x86_64"
	} else if os == "mingw" || os == "windows" {
		suffix = ".exe"
	} else {
		if arch == "x86_64" || arch == "amd64" {
			suffix = "-Linux-x86_64"
		} else if arch == "armv8*" || arch == "aarch64" {
			suffix = "-Linux-arm64"
		}
	}

	return fmt.Sprintf("https://github.com/okteto/okteto/releases/download/%s/okteto%s", oktetoVersion, suffix)
}

func DownloadOkteto() error {
	oktetoUrl := getOktetoUrl()

	if err := downloadBinary(oktetoUrl, "okteto", true); err != nil {
		return err
	}

	return nil
}
