// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
)

var oktetoUpCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long:  "",
	RunE:  runOktetoUp,
}

func init() {
	oktetoCmd.AddCommand(oktetoUpCmd)
}

func runOktetoUp(cmd *cobra.Command, args []string) error {


	return nil
}
