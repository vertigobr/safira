package get

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	c "github.com/vertigobr/safira-libs/pkg/config"
)

const helmVersion = "v3.1.2"

func getHelmUrl() string {
	arch := getARCH()
	os := getOS()

	return fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz", helmVersion, os, arch)
}

func DownloadHelm() error {
	helmUrl := getHelmUrl()
	parsedURL, _ := url.Parse(helmUrl)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	dest, err := c.CreateInBinDir()
	if err != nil {
		return err
	}

	r := ioutil.NopCloser(res.Body)
	untarErr := Untar(r, dest)
	if untarErr != nil {
		return untarErr
	}

	return nil
}
