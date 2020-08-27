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
	"gopkg.in/gookit/color.v1"
)

func PullTemplate(repo string, update, verboseFlag bool) error {
	templateFolder := "template"
	exists, err := os.Stat(templateFolder)

	if err != nil || exists != nil && update {
		if exists == nil {
			if verboseFlag {
				fmt.Printf("%s Templates not found\n", color.Blue.Text("[v]"))
			}
			fmt.Printf("%s Downloading templates\n", color.Green.Text("[↓]"))
		} else {
			if update {
				os.RemoveAll(templateFolder)
			}

			if verboseFlag {
				fmt.Printf("%s Found templates\n", color.Blue.Text("[v]"))
			}
			fmt.Printf("%s Updating templates\n", color.Green.Text("[↓]"))
		}

		dir, err := ioutil.TempDir("", "ipaasTemplates")
		if err != nil {
			return fmt.Errorf("%s Error creating temporary folder for downloading templates", color.Red.Text("[!]"))
		}
		defer os.RemoveAll(dir)

		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("%s Error downloading templates", color.Red.Text("[!]"))
		}

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("%s Error getting path to current folder", color.Red.Text("[!]"))
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
			fmt.Printf("%s Found templates\n", color.Blue.Text("[v]"))
		}

		return nil
	}
}
