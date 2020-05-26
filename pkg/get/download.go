package get

import (
	"fmt"
	"github.com/vertigobr/safira-libs/pkg/config"
	"io"
	"io/ioutil"
	"net/http"
	u "net/url"
	"os"
)

func download(url, name string, binary bool) error {
	parsedURL, _ := u.Parse(url)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	dest, err := config.CreateInBinDir()
	if err != nil {
		return err
	}

	if binary {
		// Criar arquivo
		out, err := os.Create(fmt.Sprintf("%s/%s", dest, name))
		if err != nil {
			return err
		}
		defer out.Close()

		// Escreve o corpo da resposta no arquivo
		if _, err := io.Copy(out, res.Body); err != nil {
			return err
		}

		if err := os.Chmod(fmt.Sprintf("%s/%s", dest, name), 0700); err != nil {
			return err
		}
	} else {
		r := ioutil.NopCloser(res.Body)

		if err := Untar(r, dest); err != nil {
			return err
		}
	}

	return nil
}
