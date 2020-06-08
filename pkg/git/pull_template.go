// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package git

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func PullTemplate(repo string, verboseFlag bool) error {
	exists, err := os.Stat("./template")
	if err != nil || exists == nil {
		templateFolder := "template"
		if verboseFlag {
			fmt.Println("[+] Templates não encontrados")
		}

		fmt.Println("Baixando templates...")
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
