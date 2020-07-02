// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
)

var kongpluginCmd = &cobra.Command{
	Use:   "kongplugin",
	Aliases: []string{"kp"},
	Short: "Realiza ações voltadas para os Plugins do Kong",
	Long:  "Realiza ações voltadas para os Plugins do Kong",
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(kongpluginCmd)
}
