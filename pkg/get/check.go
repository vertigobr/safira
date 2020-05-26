package get

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"os/exec"
)

var errorCheck = "não foi possível baixar o pacote "

func CheckKubectl() error {
	if exists, _ := ExistsBinary("kubectl"); !exists {
		fmt.Println("Baixando kubectl...")
		if err := DownloadKubectl(); err != nil {
			return fmt.Errorf(errorCheck + "kubectl")
		}
	}

	return nil
}

func CheckK3d() error {
	if exists, _ := ExistsBinary("k3d"); !exists {
		fmt.Println("Baixando k3d...")
		if err := DownloadK3d(); err != nil {
			return fmt.Errorf(errorCheck + "k3d")
		}
	}

	return nil
}

func CheckHelm() error {
	if exists, _ := ExistsBinary("helm"); !exists {
		fmt.Println("Baixando helm...")
		if err := DownloadHelm(); err != nil {
			return fmt.Errorf(errorCheck + "helm")
		}
	}

	return nil
}

func CheckFaasCli() error {
	if exists, _ := ExistsBinary("faas-cli"); !exists {
		fmt.Println("Baixando faas-cli...")
		if err := DownloadHelm(); err != nil {
			return fmt.Errorf(errorCheck + "faas-cli")
		}
	}

	return nil
}

func ExistsBinary(binary string) (exists bool, err error) {
	path, err := exec.LookPath(fmt.Sprintf("%sbin/%s", config.GetUserDir(), binary))
	exists = path != ""
	return
}
