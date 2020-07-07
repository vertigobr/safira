// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/docker"
	"github.com/vertigobr/safira/pkg/get"
	s "github.com/vertigobr/safira/pkg/stack"
)

var functionBuildCmd = &cobra.Command{
	Use:     "build [FUNCTION_NAME]",
	Short:   "Build Docker images from functions",
	Long:    "Build Docker images from functions",
	PreRunE: preRunFunctionBuild,
	RunE:    runFunctionBuild,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionBuildCmd)
	functionBuildCmd.Flags().BoolP("all-functions", "A", false, "build all functions")
	functionBuildCmd.Flags().Bool("no-cache", false, "do not use cache when building the image")
}

func preRunFunctionBuild(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionBuild(cmd *cobra.Command, args []string) error {
	noCacheFlag, _ := cmd.Flags().GetBool("no-cache")
	all, _ := cmd.Flags().GetBool("all-functions")

	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	if err := buildFunction(stack, args, all, noCacheFlag); err != nil {
		return err
	}

	fmt.Println("\nBuild realizado com sucesso!")

	return nil
}

func buildFunction(stack *s.Stack, args []string, allFunctions, noCacheFlag bool) error {
	buildArgsStack := stack.StackConfig.BuildArgs
	functions      := stack.Functions

	if err := get.DownloadTemplate(faasTemplateRepo, false, false); err != nil {
		return err
	}

	if allFunctions {
		for functionName, f := range functions {
			var buildArgs map[string]string

			if len(f.FunctionConfig.BuildArgs) != 0 {
				buildArgs = f.FunctionConfig.BuildArgs
			} else {
				buildArgs = buildArgsStack
			}

			err := docker.Build(f.Image, functionName, f.Handler, f.Lang, noCacheFlag, buildArgs)
			if err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			functionName := args[index]
			if checkFunctionExists(functionName, functions) {
				f := functions[functionArg]
				var buildArgs map[string]string

				if len(f.FunctionConfig.BuildArgs) != 0 {
					buildArgs = f.FunctionConfig.BuildArgs
				} else {
					buildArgs = buildArgsStack
				}

				err := docker.Build(f.Image, functionName, f.Handler, f.Lang, noCacheFlag, buildArgs)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	return nil
}
