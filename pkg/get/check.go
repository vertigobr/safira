// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"os/exec"
)

func CheckBinary(binary string, downloadIfNotExist, verboseFlag bool) (bool, error) {
	exists, _ := existsBinary(binary)
	errorCheck := "não foi possível baixar o pacote: "

	if verboseFlag {
		fmt.Println("[+] Verificando " + binary)
	}

	if !exists && downloadIfNotExist {
		fmt.Println("Baixando " + binary + "...")
		switch binary {
		case "kubectl":
			if err := DownloadKubectl(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "kubectl. Tente novamente")
			}
		case "k3d":
			if err := DownloadK3d(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "k3d. Tente novamente")
			}
		case "helm":
			if err := DownloadHelm(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "helm. Tente novamente")
			}
		case "faas-cli":
			if err := DownloadFaasCli(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "faas-cli. Tente novamente")
			}
		default:
			return false, fmt.Errorf("nome de binário inválido")
		}
	}

	return exists, nil
}

func existsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", config.GetSafiraDir(), binary))
	exists = path != ""
	return
}
