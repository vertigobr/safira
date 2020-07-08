// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	s "github.com/vertigobr/safira/pkg/stack"
)

var functionBuildPushCmd = &cobra.Command{
	Use:     "build-push [FUNCTION_NAME]",
	Short:   "Executa o build e push das imagens",
	Long:    "Executa o build e push das imagens",
	Example: `If you want to build and push a function's Docker image, run:

    $ safira function build-push function-name

or if you want to build and push the Docker image of all the functions, execute:

    $ safira function build-push -A`,
	PreRunE:  preRunFunctionBuildPush,
	RunE:     runFunctionBuildPush,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionBuildPushCmd)
	functionBuildPushCmd.Flags().BoolP("all-functions", "A", false, "pushes all Docker images from functions to the registry")
	functionBuildPushCmd.Flags().Bool("no-cache", false, "do not use cache when building the image")
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
