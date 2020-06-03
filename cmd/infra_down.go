// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Derruba uma infraestrutura de desenvolvimento provisionada anteriormente",
	Long:  "Derruba uma infraestrutura de desenvolvimento provisionada anteriormente",
	RunE:  runInfraDown,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(downCmd)
}

func runInfraDown(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	exist, err := checkInfra(verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	k3dPath := config.GetK3dPath()

	if err := deleteCluster(k3dPath, verboseFlag); err != nil {
		return errors.New("\nNenhum cluster encontrado!")
	}

	fmt.Println("\nCluster destruído com sucesso!")
	fmt.Println()
	return nil
}

func deleteCluster(k3dPath string, verboseFlag bool) error {
	taskDeleteCluster := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"delete",
			"-n", clusterName,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Destruindo cluster local...")
	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}
