// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"gopkg.in/gookit/color.v1"
)

var infraDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Take down the provisioned development infrastructure",
	Long:  "Take down the provisioned development infrastructure",
	Example: `To destroy the locally provisioned cluster, run:

    $ safira infra down`,
	RunE:                       runInfraDown,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraDownCmd)
}

func runInfraDown(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	err := checkInfra(verboseFlag)
	if err != nil {
		return err
	}

	k3dPath := config.GetK3dPath()

	if err := deleteCluster(k3dPath, verboseFlag); err != nil {
		return fmt.Errorf("%s No clusters found", color.Red.Text("[!]"))
	}

	fmt.Printf("%s Cluster destroyed successfully\n", color.Green.Text("[+]"))
	return nil
}

func deleteCluster(k3dPath string, verboseFlag bool) error {
	taskDeleteCluster := execute.Task{
		Command: k3dPath,
		Args: []string{
			"delete",
			"-n", clusterName,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Printf("%s Destroying local cluster\n", color.Green.Text("[+]"))
	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}
