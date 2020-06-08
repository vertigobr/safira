// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	s "github.com/vertigobr/safira/pkg/stack"
	"os"

	"github.com/spf13/cobra"
)

var buildPushCmd = &cobra.Command{
	Use:     "build-push [FUNCTION_NAME]",
	Short:   "Executa o build e push das imagens",
	Long:    "Executa o build e push das imagens",
	PreRunE:  preRunFunctionBuildPush,
	RunE:     runFunctionBuildPush,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(buildPushCmd)
	buildPushCmd.Flags().Bool("no-cache", false, "Do not use cache when building the image")
	buildPushCmd.Flags().BoolP("all-functions", "A", false, "Push all functions")
}

func preRunFunctionBuildPush(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionBuildPush(cmd *cobra.Command, args []string) error {
	noCacheFlag, _ := cmd.Flags().GetBool("no-cache")
	all, _ := cmd.Flags().GetBool("all-functions")

	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	if err := buildFunction(stack, args, all, noCacheFlag); err != nil {
		return err
	}

	if err := pushImage(stack, args, all); err != nil {
		return err
	}

	return nil
}
