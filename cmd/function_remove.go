// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
	"github.com/vertigobr/safira/pkg/stack"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove uma ou mais funções do cluster",
	Long:    "Remove uma ou mais funções do cluster",
	PreRunE: preRunFunctionRemove,
	RunE:    runFunctionRemove,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("all-functions", "A", false, "Deploy all functions")
}

func preRunFunctionRemove(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionRemove(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	all, _ := cmd.Flags().GetBool("all-functions")

	k8sClient, err := k8s.GetClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("cluster local não encontrado!\n")
	}

	functions, err := stack.GetAllFunctions()
	if err != nil {
		return err
	}

	if all {
		for index, _ := range functions {
			if err := removeFunction(k8sClient, index, verboseFlag); err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				if err := removeFunction(k8sClient, functionArg, verboseFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	return nil
}

func removeFunction(client *kubernetes.Clientset, functionName string, verboseFlag bool) error {
	deploymentsFunctions := client.AppsV1().Deployments(functionsNamespace)
	listFunction, _ := deploymentsFunctions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações das funções no cluster")
	}

	for _, deploy := range listFunction.Items{
		if deploy.Name == functionName {
			err := deploymentsFunctions.Delete(context.TODO(), functionName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("Function %s removida!", functionName))
			return nil
		}
	}

	return fmt.Errorf("function %s não encontrada!\n", functionName)
}
