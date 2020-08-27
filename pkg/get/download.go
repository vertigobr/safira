// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	u "net/url"
	"os"
	"os/user"
	"strconv"

	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/git"
	"gopkg.in/gookit/color.v1"
)

func downloadBinary(url, name string, binary bool) error {
	parsedURL, _ := u.Parse(url)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("%s Error when obtaining request body: %s", color.Red.Text("[!]"), url)
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
			return fmt.Errorf("%s Error when creating %s file", color.Red.Text("[!]"), name)
		}
		defer out.Close()

		// Escreve o corpo da resposta no arquivo
		if _, err := io.Copy(out, res.Body); err != nil {
			return fmt.Errorf("%s Error writing request body to file %s", color.Red.Text("[!]"), name)
		}

		if err := os.Chmod(fmt.Sprintf("%s/%s", dest, name), 0700); err != nil {
			return fmt.Errorf("%s Error when making the X file an executable - %s", color.Red.Text("[!]"), name)
		}
	} else {
		r := ioutil.NopCloser(res.Body)

		if err := Untar(r, dest); err != nil {
			return err
		}
	}

	sudoUser := os.Getenv("SUDO_USER")
	if len(sudoUser) == 0 {
		return nil
	}

	u, err := user.Lookup(sudoUser)

	Uid, err := strconv.Atoi(u.Uid)
	Gid, err := strconv.Atoi(u.Gid)

	if err := os.Chown(fmt.Sprintf("%s/%s", dest, name), Uid, Gid); err != nil {
		return fmt.Errorf("%s Error changing the owner of the root folder for the user", color.Red.Text("[!]"))
	}

	return nil
}

func DownloadTemplate(faasTemplateRepo string, update, verboseFlag bool) error {
	if err := git.PullTemplate(faasTemplateRepo, update, verboseFlag); err != nil {
		return err
	}

	return nil
}
