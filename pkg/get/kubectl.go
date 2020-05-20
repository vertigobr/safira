package get

import (
	"fmt"
	c "github.com/vertigobr/safira-libs/pkg/config"
	"io"
	"net/http"
	"net/url"
	o "os"
)

const kubectlVersion = "v1.18.0"

func getKubectlUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", kubectlVersion, os, arch)
}

func DownloadKubectl() error {
	kubectlUrl := getKubectlUrl()
	parsedURL, _ := url.Parse(kubectlUrl)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	dest, err := c.CreateInBinDir()
	if err != nil {
		return err
	}

	// Criar arquivo
	out, err := o.Create(fmt.Sprintf("%s/kubectl", dest))
	if err != nil {
		return err
	}
	defer out.Close()

	// Escreve o corpo da resposta no arquivo
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	if err := o.Chmod(fmt.Sprintf("%s/kubectl", dest), 0700); err != nil {
		return err
	}

	return nil
}
