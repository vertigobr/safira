// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/vertigobr/safira/pkg/utils"
)

func PullTemplate(repo string, update, verboseFlag bool) error {
	templateFolder := "template"
	exists, err := os.Stat(templateFolder)

	if err != nil || exists != nil && update {
		if exists == nil {
			fmt.Println("Baixando templates...")
			if verboseFlag {
				fmt.Println("[+] Templates não encontrados")
			}
		} else {
			if update {
				os.RemoveAll(templateFolder)
			}

			fmt.Println("Atualizando templates...")
			if verboseFlag {
				fmt.Println("[+] Templates encontrados")
			}
		}

		dir, err := ioutil.TempDir("", "ipaasTemplates")
		if err != nil {
			return fmt.Errorf("error ao criar pasta temporária para download dos templates")
		}
		defer os.RemoveAll(dir) // clean up

		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("error ao baixar os templates")
		}

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error ao obter pasta atual")
		}

		templateDir := filepath.Join(dir, templateFolder)
		currentDir = filepath.Join(currentDir, templateFolder)
		err = utils.Copy(templateDir, currentDir, true, true)
		if err != nil {
			return err
		}

		return nil
	} else {
		if verboseFlag {
			fmt.Println("[+] Templates encontrados")
		}

		return nil
	}
}
