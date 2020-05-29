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
				return false, fmt.Errorf(errorCheck + "kubectl")
			}
		case "k3d":
			if err := DownloadK3d(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "k3d")
			}
		case "helm":
			if err := DownloadHelm(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "helm")
			}
		case "faas-cli":
			if err := DownloadFaasCli(); err != nil {
				if verboseFlag {
					fmt.Println(err)
				}
				return false, fmt.Errorf(errorCheck + "faas-cli")
			}
		default:
			return false, fmt.Errorf("nome de binário inválido")
		}
	}

	return exists, nil
}

func existsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", config.GetUserDir(), binary))
	exists = path != ""
	return
}
