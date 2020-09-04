// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/docker"
	s "github.com/vertigobr/safira/pkg/stack"
	"gopkg.in/gookit/color.v1"
)

var functionPushCmd = &cobra.Command{
	Use:   "push [FUNCTION_NAME]",
	Short: "Pushes Docker images from the function",
	Long:  "Pushes Docker images from the function",
	Example: `If you want to push a function's Docker image, run:

    $ safira function push function-name

or if you want to push the Docker image of all the functions, execute:

    $ safira function push -A`,
	PreRunE:                    preRunFunctionPush,
	RunE:                       runFunctionPush,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionPushCmd)
	functionPushCmd.Flags().BoolP("all-functions", "A", false, "push all function Docker images")
	functionPushCmd.Flags().StringP("env", "e", "", "Set stack env file")
}

func preRunFunctionPush(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionPush(cmd *cobra.Command, args []string) error {
	allFlag, _ := cmd.Flags().GetBool("all-functions")
	envFlag, _ := cmd.Flags().GetString("env")

	stack, err := s.LoadStackFile(envFlag)
	if err != nil {
		return err
	}

	if err := pushImage(stack, args, allFlag); err != nil {
		return err
	}

	fmt.Printf("\n%s Push successfully completed\n", color.Cyan.Text("[✓]"))

	return nil
}

func pushImage(stack *s.Stack, args []string, allFunctions bool) error {
	functions := stack.Functions
	if allFunctions {
		for functionName, f := range functions {
			fmt.Printf("%s Starting push of function %s\n", color.Green.Text("[+]"), functionName)
			err := docker.Push(f.Image, f.FunctionConfig.Build.UseSha)
			if err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			functionName := args[index]
			if checkFunctionExists(functionName, functions) {
				f := functions[functionArg]

				fmt.Printf("%s Starting push of function %s\n", color.Green.Text("[+]"), functionName)
				err := docker.Push(f.Image, f.FunctionConfig.Build.UseSha)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%s Function %s does not exist", color.Red.Text("[!]"), functionArg)
			}
		}
	}

	return nil
}
