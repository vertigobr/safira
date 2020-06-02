// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lista os templates oficiais do Vertigo iPaaS",
	Long:    "Lista os templates oficiais do Vertigo iPaaS",
	RunE:    runTemplateList,
	SuggestionsMinimumDistance: 1,
}

func init() {
	templateCmd.AddCommand(listCmd)
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	exist, err := get.CheckBinary(faasBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	faasCliPath := config.GetFaasCliPath()

	if err := templateList(faasCliPath, verboseFlag); err != nil {
		return err
	}

	return nil
}

func templateList(faasCliPath string, verboseFlag bool) error {
	taskList := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"template", "store", "list", "--url", faasTemplateStoreURL,
		},
		StreamStdio:  true,
		PrintCommand: verboseFlag,
	}

	res, err := taskList.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}
