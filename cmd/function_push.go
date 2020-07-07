// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/docker"
	s "github.com/vertigobr/safira/pkg/stack"
)

var functionPushCmd = &cobra.Command{
	Use:     "push [FUNCTION_NAME]",
	Short:   "Pushes Docker images from the function",
	Long:    "Pushes Docker images from the function",
	PreRunE: preRunFunctionPush,
	RunE:    runFunctionPush,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionPushCmd)
	functionPushCmd.Flags().BoolP("all-functions", "A", false, "push all function Docker images")
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
	all, _ := cmd.Flags().GetBool("all-functions")

	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	if err := pushImage(stack, args, all); err != nil {
		return err
	}

	fmt.Println("\nPush realizado com sucesso!")

	return nil
}

func pushImage(stack *s.Stack, args []string, allFunctions bool) error {
	functions := stack.Functions
	if allFunctions {
		for _, f := range functions {
			err := docker.Push(f.Image)
			if err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			functionName := args[index]
			if checkFunctionExists(functionName, functions) {
				f := functions[functionArg]

				err := docker.Push(f.Image)
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
