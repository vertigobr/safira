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
	"gopkg.in/gookit/color.v1"
)

var functionBuildCmd = &cobra.Command{
	Use:   "build [FUNCTION_NAME]",
	Short: "Build Docker images from functions",
	Long:  "Build Docker images from functions",
	Example: `If you want to build a function's Docker image, run:

    $ safira function build function-name

or if you want to build the Docker image of all the functions, execute:

    $ safira function build -A`,
	PreRunE:                    preRunFunctionBuild,
	RunE:                       runFunctionBuild,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionBuildCmd)
	functionBuildCmd.Flags().Bool("no-cache", false, "do not use cache when building the image")
	functionBuildCmd.Flags().BoolP("all-functions", "A", false, "build all functions")
	functionBuildCmd.Flags().Bool("update-template", false, "update template")
	functionBuildCmd.Flags().StringP("env", "e", "", "set stack env file")
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

	fmt.Printf("\n%s Build successfully completed\n", color.Cyan.Text("[✓]"))

	return nil
}

func buildFunction(stack *s.Stack, args []string, allFunctions, updateTemplateFlag, noCacheFlag, verboseFlag bool) (bool, error) {
	if stack.StackConfig.Build.Enabled != nil {
		if !*stack.StackConfig.Build.Enabled {
			fmt.Printf("%s All functions skipped the build process\n", color.Yellow.Text("[*]"))
			return true, nil
		}
	}

	buildArgsStack := stack.StackConfig.Build.Args
	functions := stack.Functions

	if err := get.DownloadTemplate(faasTemplateRepo, updateTemplateFlag, false); err != nil {
		return false, err
	}

	if allFunctions {
		for functionName, f := range functions {
			if f.FunctionConfig.Build.Enabled != nil {
				if !*f.FunctionConfig.Build.Enabled {
					fmt.Printf("%s Function %s skipped the build process\n", color.Yellow.Text("[*]"), functionName)
					continue
				}
			}

			var buildArgs map[string]string

			if len(f.FunctionConfig.Build.Args) != 0 {
				buildArgs = f.FunctionConfig.Build.Args
			} else {
				buildArgs = buildArgsStack
			}

			fmt.Printf("%s Starting build of function %s\n", color.Green.Text("[+]"), functionName)
			useSha := f.FunctionConfig.Build.UseSha || stack.StackConfig.Build.UseSha
			err := docker.Build(f.Image, functionName, f.Handler, f.Lang, useSha, noCacheFlag, buildArgs, verboseFlag)
			if err != nil {
				return false, err
			}
		}
	} else {
		for index, functionArg := range args {
			functionName := args[index]
			if checkFunctionExists(functionName, functions) {
				f := functions[functionArg]

				if f.FunctionConfig.Build.Enabled != nil {
					if !*f.FunctionConfig.Build.Enabled {
						fmt.Printf("%s Function %s skipped the build process\n", color.Yellow.Text("[*]"), functionName)
						continue
					}
				}

				var buildArgs map[string]string

				if len(f.FunctionConfig.Build.Args) != 0 {
					buildArgs = f.FunctionConfig.Build.Args
				} else {
					buildArgs = buildArgsStack
				}

				fmt.Printf("%s Starting build of function %s\n", color.Green.Text("[+]"), functionName)
				useSha := f.FunctionConfig.Build.UseSha || stack.StackConfig.Build.UseSha
				err := docker.Build(f.Image, functionName, f.Handler, f.Lang, useSha, noCacheFlag, buildArgs, verboseFlag)
				if err != nil {
					return false, err
				}
			} else {
				return false, fmt.Errorf("%s Function name %s is invalid", color.Red.Text("[!]"), functionArg)
			}
		}
	}

	return false, nil
}
