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
	msgError := "failed to update safira"

	if verboseFlag {
		fmt.Println("[+] Creating temp folder")
	}

	tmpFile, err := ioutil.TempFile("", "safira-upgrade.*.sh")
	if err != nil {
		return fmt.Errorf(msgError)
	}
	defer os.RemoveAll(tmpFile.Name())

	if verboseFlag {
		fmt.Println("[+] Downloading script")
	}

	response, err := http.Get("https://raw.githubusercontent.com/vertigobr/safira/master/install.sh")
	if err != nil {
		return fmt.Errorf(msgError)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if _, err := tmpFile.Write(body); err != nil {
		return fmt.Errorf(msgError)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf(msgError)
	}

	_ = tmpFile.Chmod(0700)

	taskUpgrade := execute.Task{
		Command: "bash",
		Args: []string{
			tmpFile.Name(),
		},
		StreamStdio: true,
	}

	fmt.Println("Upgrade Safira...")
	resUpgrade, err := taskUpgrade.Execute()
	if err != nil {
		return err
	}

	if resUpgrade.ExitCode != 0 {
		return fmt.Errorf(resUpgrade.Stderr)
	}

	return nil
}
