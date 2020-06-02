// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
)

var functionCmd = &cobra.Command{
	Use:     "function",
	Aliases: []string{"func"},
	Short:   "Permite gerenciamento das functions",
	Long:    "Permite gerenciamento das functions",
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(functionCmd)
}
