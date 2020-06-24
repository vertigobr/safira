// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
	"github.com/vertigobr/safira/pkg/stack"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var functionRemoveCmd = &cobra.Command{
	Use:     "remove [FUNCTION_NAME]",
	Aliases: []string{"rm"},
	Short:   "Remove uma ou mais funções do cluster",
	Long:    "Remove uma ou mais funções do cluster",
	PreRunE: preRunFunctionRemove,
	RunE:    runFunctionRemove,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionRemoveCmd)
	functionRemoveCmd.Flags().BoolP("all-functions", "A", false, "Deploy all functions")
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
			if err := removeDeployment(k8sClient, index, functionsNamespace, "Function", verboseFlag); err != nil {
				return err
			}

			if swaggerFile := checkSwaggerFileExist(functions[index].Handler); len(swaggerFile) > 1 {
				if err := removeDeployment(k8sClient, index + "-swagger-ui", "default", "Swagger", verboseFlag); err != nil {
					return err
				}
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				if err := removeDeployment(k8sClient, functionArg, functionsNamespace, "Function", verboseFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	return nil
}

func removeDeployment(client *kubernetes.Clientset, deployName, namespace, title string, verboseFlag bool) error {
	deploymentsFunctions := client.AppsV1().Deployments(namespace)
	listFunction, _ := deploymentsFunctions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Obtendo informações das funções no cluster")
	}

	for _, deploy := range listFunction.Items{
		if deploy.Name == deployName {
			err := deploymentsFunctions.Delete(context.TODO(), deployName, v1.DeleteOptions{})
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s %s removida!", title, deployName))
			return nil
		}
	}

	return nil // fmt.Errorf("%s %s não encontrada!\n", title, deployName)
}
