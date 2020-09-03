// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
)

var functionLogCmd = &cobra.Command{
	Use:   "log [FUNCTION_NAME]",
	Short: "Output a function log",
	Long:  "Output a function log",
	Example: `To view a function log, run:

    $ safira function log function-name`,
	PreRunE:                    preRunFunctionLog,
	RunE:                       runFunctionLog,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionLogCmd)
	functionLogCmd.Flags().String("kubeconfig", kubeconfigPath, "set kubeconfig to deploy")
	functionLogCmd.Flags().StringP("namespace", "n", functionsNamespace, "set namespace to deploy")
	functionLogCmd.Flags().StringP("output", "o", "", "set output file")
}

func preRunFunctionLog(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionLog(cmd *cobra.Command, args []string) error {
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	namespaceFlag, _ := cmd.Flags().GetString("namespace")
	outputFlag, _ := cmd.Flags().GetString("output")

	if err := k8s.OutputFunctionLog(args[0], kubeconfigFlag, namespaceFlag, outputFlag); err != nil {
		return err
	}

	return nil
}
