// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
	"github.com/vertigobr/safira/pkg/stack"
)

var functionUndeployCmd = &cobra.Command{
	Use:     "undeploy [FUNCTION_NAME]",
	Short:   "Remove a function from the cluster",
	Long:    "Remove a function from the cluster",
	Example: `To remove the function from a project, run:

    $ safira function undeploy function-name

or if you want to remove all functions from a project, execute:

    $ safira function undeploy -A`,
	PreRunE: preRunFunctionUndeploy,
	RunE:    runFunctionUndeploy,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionUndeployCmd)
	functionUndeployCmd.Flags().BoolP("all-functions", "A", false, "deploy all functions")
	functionUndeployCmd.Flags().String("kubeconfig", kubeconfigPath, "set kubeconfig to remove function")
	functionUndeployCmd.Flags().StringP("namespace", "n", functionsNamespace, "set namespace to undeploy")
}

func preRunFunctionUndeploy(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionUndeploy(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	all, _ := cmd.Flags().GetBool("all-functions")
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	namespaceFlag, _ := cmd.Flags().GetString("namespace")

	functions, err := stack.GetAllFunctions()
	if err != nil {
		return err
	}

	if all {
		for index, _ := range functions {
			if err := removeDeploy(index, namespaceFlag, kubeconfigFlag, "Function", verboseFlag); err != nil {
				return err
			}

			if swaggerFile := checkSwaggerFileExist(functions[index].Handler); len(swaggerFile) > 1 {
				if err := removeDeploy(index + "-swagger-ui", kubeconfigFlag, "default", "Swagger", verboseFlag); err != nil {
					return err
				}
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				if err := removeDeploy(functionArg, namespaceFlag, kubeconfigFlag, "Function", verboseFlag); err != nil {
					return err
				}

				if swaggerFile := checkSwaggerFileExist(functions[functionArg].Handler); len(swaggerFile) > 1 {
					if err := removeDeploy(functionArg + "-swagger-ui", kubeconfigFlag, "default", "Swagger", verboseFlag); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	return nil
}

func removeDeploy(name, namespace, kubeconfigFlag, title string, verboseFlag bool) error {
	k8sClient, err := k8s.GetClient(kubeconfigFlag)
	if err != nil {
		return fmt.Errorf("cluster não encontrado!\n")
	}

	if err := k8s.RemoveDeployment(k8sClient, name, namespace, title, verboseFlag); err != nil {
		return err
	}

	if err := k8s.RemoveService(k8sClient, name, namespace, title, verboseFlag); err != nil {
		return err
	}

	if err := k8s.RemoveIngress(k8sClient, name, namespace, title, verboseFlag); err != nil {
		return err
	}

	return nil
}
