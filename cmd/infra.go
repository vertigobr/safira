// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
)

var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Responsável por gerenciar a infraestrutura",
	Long:  "Responsável por gerenciar a infraestrutura em ambiente local", //  ou em nuvem
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(infraCmd)

	infraCmd.PersistentFlags().String(
		"env",
		"local",
		"Recebe o ambiente aonde será provisionado o cluster Kubernetes.",
	)
}

func checkInfra(verboseFlag bool) (bool, error) {
	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return exist, err
	}

	if !exist {
		return exist, nil
	}

	exist, err = get.CheckBinary(k3dBinaryName, false, verboseFlag)
	if err != nil {
		return exist, err
	}

	if !exist {
		return exist, nil
	}

	exist, err = get.CheckBinary(helmBinaryName, false, verboseFlag)
	if err != nil {
		return exist, err
	}

	if !exist {
		return exist, nil
	}

	return true, nil
}
