// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/execute"
	"gopkg.in/gookit/color.v1"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "CLI updates itself",
	Long:  "CLI updates itself",
	Example: `To upgrade CLI, run:

    $ safira upgrade`,
	RunE: runUpgrade,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringP("version", "v", "latest", "CLI download version")
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	if err := upgrade(verboseFlag); err != nil {
		return err
	}

	return nil
}

func upgrade(verboseFlag bool) error {
	if verboseFlag {
		fmt.Printf("%s Creating temp folder\n", color.Blue.Text("[v]"))
	}

	tmpFile, err := ioutil.TempFile("", "safira-upgrade.*.sh")
	if err != nil {
		return fmt.Errorf("%s Error creating temporary file", color.Red.Text("[!]"))
	}
	defer os.RemoveAll(tmpFile.Name())

	if verboseFlag {
		fmt.Printf("%s Downloading script\n", color.Blue.Text("[v]"))
	}

	response, err := http.Get("https://raw.githubusercontent.com/vertigobr/safira/master/install.sh")
	if err != nil {
		return fmt.Errorf("%s Error while downloading the update script", color.Red.Text("[!]"))
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if _, err := tmpFile.Write(body); err != nil {
		return fmt.Errorf("%s Error saving the update script", color.Red.Text("[!]"))
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("%s Error when closing the update script", color.Red.Text("[!]"))
	}

	_ = tmpFile.Chmod(0700)

	taskUpgrade := execute.Task{
		Command: "bash",
		Args: []string{
			tmpFile.Name(),
		},
		StreamStdio: true,
	}

	fmt.Printf("%s Upgrade Safira\n", color.Green.Text("[+]"))
	resUpgrade, err := taskUpgrade.Execute()
	if err != nil {
		return err
	}

	if resUpgrade.ExitCode != 0 {
		return fmt.Errorf(resUpgrade.Stderr)
	}

	return nil
}
