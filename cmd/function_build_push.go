// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	s "github.com/vertigobr/safira/pkg/stack"
	"gopkg.in/gookit/color.v1"
)

var functionBuildPushCmd = &cobra.Command{
	Use:     "build-push [FUNCTION_NAME]",
	Aliases: []string{"bp"},
	Short:   "Executa o build e push das imagens",
	Long:    "Executa o build e push das imagens",
	Example: `If you want to build and push a function's Docker image, run:

    $ safira function build-push function-name

or if you want to build and push the Docker image of all the functions, execute:

    $ safira function build-push -A`,
	PreRunE:                    preRunFunctionBuildPush,
	RunE:                       runFunctionBuildPush,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionBuildPushCmd)
	functionBuildPushCmd.Flags().Bool("no-cache", false, "do not use cache when building the image")
	functionBuildPushCmd.Flags().BoolP("all-functions", "A", false, "pushes all Docker images from functions to the registry")
	functionBuildPushCmd.Flags().Bool("update-template", false, "update template")
	functionBuildPushCmd.Flags().StringP("env", "e", "", "Set stack env file")
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
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	noCacheFlag, _ := cmd.Flags().GetBool("no-cache")
	allFlag, _ := cmd.Flags().GetBool("all-functions")
	envFlag, _ := cmd.Flags().GetString("env")
	updateTemplateFlag, _ := cmd.Flags().GetBool("update-template")

	stack, err := s.LoadStackFile(envFlag)
	if err != nil {
		return err
	}

	if skipped, err := buildFunction(stack, args, allFlag, updateTemplateFlag, noCacheFlag, verboseFlag); err != nil {
		return err
	} else if skipped {
		os.Exit(0)
	}

	if err := pushImage(stack, args, allFlag); err != nil {
		return err
	}

	fmt.Printf("\n%s Build and Push successfully completed\n", color.Cyan.Text("[✓]"))

	return nil
}
