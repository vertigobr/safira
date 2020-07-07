// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage the use of templates",
	Long:  "Manage the use of templates",
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
