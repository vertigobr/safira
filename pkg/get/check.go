package get

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"os/exec"
)

//func CheckKubectl(downloadIfNotExist bool) error {
//	if exists, _ := existsBinary("kubectl"); !exists {
//			fmt.Println("Baixando kubectl...")
//		if err := DownloadKubectl(); err != nil {
//			return fmt.Errorf(errorCheck + "kubectl")
//		}
//	}
//
//	return nil
//}
//
//func CheckK3d(downloadIfNotExist bool) error {
//	if exists, _ := existsBinary("k3d"); !exists {
//		fmt.Println("Baixando k3d...")
//		if err := DownloadK3d(); err != nil {
//			return fmt.Errorf(errorCheck + "k3d")
//		}
//	}
//
//	return nil
//}
//
//func CheckHelm(downloadIfNotExist bool) error {
//	if exists, _ := existsBinary("helm"); !exists && downloadIfNotExist {
//		fmt.Println("Baixando helm...")
//		if err := DownloadHelm(); err != nil {
//			return fmt.Errorf(errorCheck + "helm")
//		}
//	}
//
//	return nil
//}
//
//func CheckFaasCli(downloadIfNotExist bool) error {
//	if exists, _ := existsBinary("faas-cli"); !exists && downloadIfNotExist {
//		fmt.Println("Baixando faas-cli...")
//		if err := DownloadFaasCli(); err != nil {
//			return fmt.Errorf(errorCheck + "faas-cli")
//		}
//	}
//
//	return nil
//}

func CheckBinary(binary string, downloadIfNotExist bool) (bool, error) {
	fmt.Println("Verificando dependências...")
	exists, _ := existsBinary(binary)
	errorCheck := "não foi possível baixar o pacote: "
	if !exists && downloadIfNotExist {
		fmt.Println("Baixando " + binary + "...")
		switch binary {
		case "kubectl":
			if err := DownloadKubectl(); err != nil {
				return false, fmt.Errorf(errorCheck + "kubectl")
			}
		case "k3d":
			if err := DownloadK3d(); err != nil {
				return false, fmt.Errorf(errorCheck + "k3d")
			}
		case "helm":
			if err := DownloadHelm(); err != nil {
				return false, fmt.Errorf(errorCheck + "helm")
			}
		case "faas-cli":
			if err := DownloadFaasCli(); err != nil {
				return false, fmt.Errorf(errorCheck + "faas-cli")
			}
		}
	}

	return exists, nil
}

func existsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", config.GetUserDir(), binary))
	exists = path != ""
	return
}
