// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/vertigobr/safira/pkg/k8s"
	"os"

	"github.com/spf13/cobra"
)

var functionLogCmd = &cobra.Command{
	Use:     "log [FUNCTION_NAME]",
	Short:   "Imprime o log de uma função",
	Long:    "Imprime o log de uma função",
	PreRunE: preRunFunctionLog,
	RunE:    runFunctionLog,
}

func init() {
	functionCmd.AddCommand(functionLogCmd)
	functionLogCmd.Flags().String("kubeconfig", kubeconfigPath, "Set kubeconfig to deploy")
	functionLogCmd.Flags().StringP("namespace", "n", functionsNamespace, "Set namespace to deploy")
	functionLogCmd.Flags().StringP("output", "o", "", "Set output file")
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
