// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/stack"
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

func checkFunctionExists(functionName string, functions map[string]stack.Function) bool {
	f := functions[functionName]
	return len(f.Handler) > 1 && len(f.Image) > 1
}
