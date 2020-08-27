// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
	"os/exec"

	"github.com/vertigobr/safira/pkg/config"
	"gopkg.in/gookit/color.v1"
)

func CheckBinary(binary string, downloadIfNotExist, verboseFlag bool) (bool, error) {
	exists, _ := existsBinary(binary)

	if verboseFlag {
		fmt.Printf("%s Checking the binary: %s\n", color.Blue.Text("[v]"), binary)
	}

	if !exists && downloadIfNotExist {
		fmt.Printf("%s Downloading %s\n", color.Green.Text("[↓]"), binary)
		switch binary {
		case "kubectl":
			if err := DownloadKubectl(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, downloadMessageError("kubectl")
			}
		case "k3d":
			if err := DownloadK3d(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, downloadMessageError("k3d")
			}
		case "helm":
			if err := DownloadHelm(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, downloadMessageError("helm")
			}
		case "faas-cli":
			if err := DownloadFaasCli(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, downloadMessageError("faas-cli")
			}
		case "okteto":
			if err := DownloadOkteto(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, downloadMessageError("okteto")
			}
		default:
			return false, fmt.Errorf("%s Invalid binary name", color.Red.Text("[!]"))
		}
	}

	return exists, nil
}

func existsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", config.GetSafiraDir(), binary))
	exists = path != ""
	return
}

func downloadMessageError(binary string) error {
	return fmt.Errorf("%s Unable to download %s, try again", color.Red.Text("[!]"), binary)
}
