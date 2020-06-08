// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var kubeconfigCmd = &cobra.Command{
	Use:   "get-kubeconfig",
	Short: "Imprime o caminho do kubeconfig",
	Long:  "Imprime o caminho do kubeconfig",
	RunE:  runGetKubeconfig,
}

func init() {
	infraCmd.AddCommand(kubeconfigCmd)
}

func runGetKubeconfig(cmd *cobra.Command, args []string) error {
	fmt.Println(kubeconfigPath)

	return nil
}
