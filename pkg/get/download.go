// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/git"
	"io"
	"io/ioutil"
	"net/http"
	u "net/url"
	"os"
	"os/user"
	"strconv"
)

func downloadBinary(url, name string, binary bool) error {
	parsedURL, _ := u.Parse(url)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("error ao obter conteúdo da requisição: %s", err.Error())
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
			return fmt.Errorf("error ao criar arquivo %s: %s", name, err.Error())
		}
		defer out.Close()

		// Escreve o corpo da resposta no arquivo
		if _, err := io.Copy(out, res.Body); err != nil {
			return fmt.Errorf("error ao escrever conteúdo da requisição no arquivo %s: %s", name, err.Error())
		}

		if err := os.Chmod(fmt.Sprintf("%s/%s", dest, name), 0700); err != nil {
			return fmt.Errorf("error ao tornar o arquivo executavél - %s: %s", name, err.Error())
		}
	} else {
		r := ioutil.NopCloser(res.Body)

		if err := Untar(r, dest); err != nil {
			return err
		}
	}

	sudoUser := os.Getenv("SUDO_USER")

	u, err := user.Lookup(sudoUser)

	Uid, err := strconv.Atoi(u.Uid)
	Gid, err := strconv.Atoi(u.Gid)

	if err := os.Chown(fmt.Sprintf("%s/%s", dest, name), Uid, Gid); err != nil {
		return fmt.Errorf("error ao mudar o dono da pasta de root para o usuário: %s", err.Error())
	}

	return nil
}

func DownloadTemplate(faasTemplateRepo string, verboseFlag bool) error {
	if err := git.PullTemplate(faasTemplateRepo, verboseFlag); err != nil {
		return err
	}

	return nil
}
