// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
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
	Use:     "deploy -f YAML_FILE",
	Short:   "Executa deploy das funções",
	Long:    "Executa deploy das funções",
	RunE:    runFunctionDeploy,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(deployCmd)
	deployCmd.Flags().Bool("update", false, "Force the deploy to pull a new image. (Default: false)")
	deployCmd.Flags().String("kubeconfig", "", "Set kubeconfig to deploy")
}

func runFunctionDeploy(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	updateFlag, _ := cmd.Flags().GetBool("update")
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()

	if err := checkDeployFiles(); err!= nil {
		return err
	}

	if err := functionDeploy(kubectlPath, kubeconfigFlag, verboseFlag, updateFlag); err != nil {
		return err
	}

	fmt.Println("\nDeploy realizado com sucesso!")

	return nil
}

func functionDeploy(kubectlPath, kubeconfigFlag string, verboseFlag, updateFlag bool) error {
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

	hasFunction, err := deploy.CheckFunction(clusterName)
	if err != nil {
		return err
	}

	projectName, err := deploy.GetProjectName()
	if err != nil {
		return err
	}

	if hasFunction && updateFlag {
		taskRemoveFunction := execute.Task{
			Command:     kubectlPath,
			Args:        []string{
				"rollout", "restart", "deployments", projectName,
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

func checkDeployFiles() error {
	deployFolder := "./deploy"

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	if _, err := os.Stat(".env"); err != nil {
		return fmt.Errorf("arquivo .env não encontrado")
	}

	if err := deploy.CreateYamlFunction(deployFolder + "/function.yml"); err != nil {
		return err
	}

	if err := deploy.CreateYamlIngress(deployFolder + "/ingress.yml"); err != nil {
		return err
	}

	if err := deploy.CreateYamlKongPlugin(deployFolder + "/kongplugin.yml"); err != nil {
		return err
	}

	if err := deploy.CreateYamlService(deployFolder + "/" + "service.yml"); err != nil {
		return err
	}

	return nil
}
