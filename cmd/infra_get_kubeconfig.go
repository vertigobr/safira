// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infraKubeconfigCmd = &cobra.Command{
	Use:   "get-kubeconfig",
	Short: "Output kubeconfig path",
	Long:  "Output kubeconfig path",
	Example: `To output kubeconfig path, run:

    $ safira infra get-kubeconfig`,
	RunE:                       runGetKubeconfig,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraKubeconfigCmd)
}

func runGetKubeconfig(_ *cobra.Command, _ []string) error {
	fmt.Println(kubeconfigPath)

	return nil
}
