// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/get"
	"github.com/vertigobr/safira/pkg/stack"
	"os"

	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/deploy"
	"github.com/vertigobr/safira/pkg/execute"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy [FUNCTION_NAME]",
	Short:   "Executa deploy das funções",
	Long:    "Executa deploy das funções",
	PreRunE: preRunFunctionDeploy,
	RunE:    runFunctionDeploy,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(deployCmd)
	deployCmd.Flags().Bool("update", false, "Force the deploy to pull a new image (Default: false)")
	deployCmd.Flags().String("kubeconfig", "", "Set kubeconfig to deploy")
	deployCmd.Flags().BoolP("all-functions", "A", false, "Deploy all functions")
}

func preRunFunctionDeploy(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all-functions")
	if len(args) < 1 && !all {
		_ = cmd.Help()
		os.Exit(0)
	}

	return nil
}

func runFunctionDeploy(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	updateFlag, _ := cmd.Flags().GetBool("update")
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	all, _ := cmd.Flags().GetBool("all-functions")

	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()
	functions, err := stack.GetAllFunctions()
	if err != nil {
		return err
	}

	if all {
		for index, _ := range functions {
			if err := checkDeployFiles(index); err!= nil {
				return err
			}

			if err := functionDeploy(kubectlPath, kubeconfigFlag, index, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				if err := checkDeployFiles(functionArg); err!= nil {
					return err
				}

				if err := functionDeploy(kubectlPath, kubeconfigFlag, functionArg, verboseFlag, updateFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	fmt.Println("\nDeploy realizado com sucesso!")

	return nil
}

func functionDeploy(kubectlPath, kubeconfigFlag, functionName string, verboseFlag, updateFlag bool) error {
	var kubeconfig string
	if len(kubeconfigFlag) > 0 {
		kubeconfig = kubeconfigFlag
	} else {
		err := config.SetKubeconfig(clusterName)
		if err != nil {
			return err
		}

		kubeconfig = config.GetKubeconfig()
	}

	hasFunction, err := deploy.CheckFunction(clusterName, functionName, functionsNamespace)
	if err != nil {
		return err
	}

	if hasFunction && updateFlag {
		taskRemoveFunction := execute.Task{
			Command:     kubectlPath,
			Args:        []string{
				"rollout", "restart", "deployments", functionName,
				"-n", functionsNamespace,
				"--kubeconfig", kubeconfig,
			},
			StreamStdio:  verboseFlag,
			PrintCommand: verboseFlag,
		}

		if verboseFlag {
			fmt.Printf("[+] Reiniciando a função " + functionName)
		}

		res, err := taskRemoveFunction.Execute()
		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}
	}

	deployFolder := fmt.Sprintf("deploy/%s/", functionName)
	taskFunctionDeploy := execute.Task{
		Command:     kubectlPath,
		Args:        []string{
			"apply", "--wait",
			"--kubeconfig", kubeconfig,
			"-f", deployFolder,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Executando deploy da função " + functionName + "...")
	res, err := taskFunctionDeploy.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}

func checkDeployFiles(functionName string) error {
	deployFolder := fmt.Sprintf("./deploy/%s/", functionName)
	functionYaml := deployFolder + "function.yml"
	ingressYaml  := deployFolder + "ingress.yml"
	serviceYaml  := deployFolder + "service.yml"

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	if err := deploy.CreateYamlFunction(functionYaml, functionName, functionsNamespace); err != nil {
		return err
	}

	if err := deploy.CreateYamlIngress(ingressYaml, functionName); err != nil {
		return err
	}

	if err := deploy.CreateYamlService(serviceYaml, functionName); err != nil {
		return err
	}

	return nil
}
