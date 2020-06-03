// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/get"
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
	deployCmd.Flags().Bool("update", false, "Force the deploy to pull a new image. (Default: false)")
	deployCmd.Flags().String("kubeconfig", "", "Set kubeconfig to deploy")
}

func preRunFunctionDeploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("nome da função é obrigatório")
	}

	return nil
}

func runFunctionDeploy(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	updateFlag, _ := cmd.Flags().GetBool("update")
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	functionName := args[0]

	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()

	if err := checkDeployFiles(functionName); err!= nil {
		return err
	}

	if err := functionDeploy(kubectlPath, kubeconfigFlag, functionName, verboseFlag, updateFlag); err != nil {
		return err
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

	hasFunction, err := deploy.CheckFunction(clusterName, functionName)
	if err != nil {
		return err
	}

	if hasFunction && updateFlag {
		taskRemoveFunction := execute.Task{
			Command:     kubectlPath,
			Args:        []string{
				"rollout", "restart", "deployments", functionName,
				"-n", deploy.GetNamespaceFunction(),
				"--kubeconfig", kubeconfig,
			},
			StreamStdio:  verboseFlag,
			PrintCommand: verboseFlag,
		}

		if verboseFlag {
			fmt.Printf("[+] Reiniciando as função")
		}

		res, err := taskRemoveFunction.Execute()
		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}
	}

	taskFunctionDeploy := execute.Task{
		Command:     kubectlPath,
		Args:        []string{
			"apply", "--wait",
			"--kubeconfig", kubeconfig,
			"-f", "deploy/",
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Executando deploy da função...")
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
	deployFolder   := "./deploy"
	functionYaml   := deployFolder + "/" + functionName + "-function.yml"
	ingressYaml    := deployFolder + "/" + functionName + "-ingress.yml"
	serviceYaml    := deployFolder + "/" + functionName + "-service.yml"

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	if err := deploy.CreateYamlFunction(functionYaml, functionName); err != nil {
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
