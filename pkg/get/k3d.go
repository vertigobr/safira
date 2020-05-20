package get

import (
	"fmt"
	c "github.com/vertigobr/safira-libs/pkg/config"
	"io"
	"net/http"
	"net/url"
	o "os"
)

const k3dVersion = "v1.7.0"

func getK3dUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://github.com/rancher/k3d/releases/download/%s/k3d-%s-%s", k3dVersion, os, arch)
}

func DownloadK3d() error {
	k3dUrl := getK3dUrl()
	parsedURL, _ := url.Parse(k3dUrl)

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
	out, err := o.Create(fmt.Sprintf("%s/k3d", dest))
	if err != nil {
		return err
	}
	defer out.Close()

	// Escreve o corpo da resposta no arquivo
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	if err := o.Chmod(fmt.Sprintf("%s/k3d", dest), 0700); err != nil {
		return err
	}

	return nil
}
