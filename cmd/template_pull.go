// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
)

var templatePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull out the official Vertigo iPaaS templates",
	Long:  "Pull out the official Vertigo iPaaS templates",
	Example: `To pull the official templates, run:

    $ safira template pull`,
	RunE:  runTemplatePull,
}

func init() {
	templateCmd.AddCommand(templatePullCmd)
}

func runTemplatePull(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	if err := get.DownloadTemplate(faasTemplateRepo, true, verboseFlag); err != nil {
		return err
	}

	return nil
}
