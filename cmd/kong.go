// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
)

var kongCmd = &cobra.Command{
	Use:   "kong",
	Short: "Performs actions focused on Kong",
	Long:  "Performs actions focused on Kong",
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(kongCmd)
}
