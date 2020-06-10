// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/get"
	s "github.com/vertigobr/safira/pkg/stack"
	"os"

	"github.com/vertigobr/safira/pkg/config"
	d "github.com/vertigobr/safira/pkg/deploy"
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
	deployCmd.Flags().String("hostname", "", "Set hostname to deploy")
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
	hostnameFlag, _ := cmd.Flags().GetString("hostname")
	all, _ := cmd.Flags().GetBool("all-functions")

	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()
	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	functions := stack.Functions
	if all {
		for index, _ := range functions {
			if err := checkDeployFiles(index, hostnameFlag); err!= nil {
				return err
			}

			deployFolder := fmt.Sprintf("deploy/%s/", index)
			if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, index, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				if err := checkDeployFiles(functionArg, hostnameFlag); err!= nil {
					return err
				}

				deployFolder := fmt.Sprintf("deploy/%s/", functionArg)
				if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, functionArg, verboseFlag, updateFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	if len(stack.Custom) > 0 {
		for _, path := range stack.Custom {
			if err := deploy(kubectlPath, kubeconfigFlag, path, "", verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	}

	fmt.Println("\nDeploy realizado com sucesso!")

	return nil
}

func deploy(kubectlPath, kubeconfigFlag, deployFolder, functionName string, verboseFlag, updateFlag bool) error {
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

	if len(functionName) > 0 {
		hasFunction, err := d.CheckFunction(clusterName, functionName, functionsNamespace)
		if err != nil {
			return err
		}

		if hasFunction && updateFlag {
			if err := rolloutFunction(kubectlPath, kubeconfig, functionName, verboseFlag); err != nil {
				return err
			}
		}
	}

	taskDeploy := execute.Task{
		Command:     kubectlPath,
		Args:        []string{
			"apply", "--wait",
			"--kubeconfig", kubeconfig,
			"-f", deployFolder,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	if len(functionName) > 0 {
		fmt.Println("Executando deploy da função " + functionName + "...")
	} else {
		fmt.Println("Executando deploy de arquivos customizados, " + deployFolder)
	}

	res, err := taskDeploy.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}

func checkDeployFiles(functionName, hostnameFlag string) error {
	deployFolder := fmt.Sprintf("./deploy/%s/", functionName)
	functionYaml := deployFolder + "function.yml"
	ingressYaml  := deployFolder + "ingress.yml"
	serviceYaml  := deployFolder + "service.yml"

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	if err := d.CreateYamlFunction(functionYaml, functionName, functionsNamespace); err != nil {
		return err
	}

	if err := d.CreateYamlIngress(ingressYaml, functionName, hostnameFlag); err != nil {
		return err
	}

	if err := d.CreateYamlService(serviceYaml, functionName, hostnameFlag); err != nil {
		return err
	}

	return nil
}

func rolloutFunction(kubectlPath, kubeconfig, functionName string, verboseFlag bool) error {
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
		fmt.Println("[+] Reiniciando a função " + functionName)
	}

	res, err := taskRemoveFunction.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}
