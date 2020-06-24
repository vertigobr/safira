// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	d "github.com/vertigobr/safira/pkg/deploy"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	s "github.com/vertigobr/safira/pkg/stack"
)

var functionDeployCmd = &cobra.Command{
	Use:     "deploy [FUNCTION_NAME]",
	Short:   "Executa deploy das funções",
	Long:    "Executa deploy das funções",
	PreRunE: preRunFunctionDeploy,
	RunE:    runFunctionDeploy,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionDeployCmd)
	functionDeployCmd.Flags().Bool("update", false, "Force the deploy to pull a new image (Default: false)")
	functionDeployCmd.Flags().String("kubeconfig", "", "Set kubeconfig to deploy")
	functionDeployCmd.Flags().String("hostname", "", "Set hostname to deploy")
	functionDeployCmd.Flags().BoolP("all-functions", "A", false, "Deploy all functions")
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
			swaggerFile := checkSwaggerFileExist(functions[index].Handler)
			if err := checkDeployFiles(index, functions[index].Handler, hostnameFlag, swaggerFile); err!= nil {
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
				swaggerFile := checkSwaggerFileExist(functions[functionArg].Handler)

				if err := checkDeployFiles(functionArg, functions[functionArg].Handler, hostnameFlag, swaggerFile); err!= nil {
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

func checkDeployFiles(functionName, functionHandler, hostnameFlag , swaggerFile string) error {
	deployFolder     := fmt.Sprintf("./deploy/%s/", functionName)
	functionYamlName := deployFolder + "function.yml"
	ingressYamlName  := deployFolder + "ingress.yml"
	serviceYamlName  := deployFolder + "service.yml"
	functionPath     := fmt.Sprintf("/function/%s", functionName)

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	var functionYaml d.K8sYaml
	if err := functionYaml.MountFunction(functionName, functionsNamespace); err != nil {
		return err
	}

	if err := functionYaml.CreateYamlFile(functionYamlName); err != nil {
		return err
	}

	var ingressYaml d.K8sYaml
	if err := ingressYaml.MountIngress(functionName, functionName, functionPath, hostnameFlag); err != nil {
		return err
	}

	if err := ingressYaml.CreateYamlFile(ingressYamlName); err != nil {
		return err
	}

	var serviceYaml d.K8sYaml
	if err := serviceYaml.MountService(functionName, hostnameFlag, true); err != nil {
		return err
	}

	if err := serviceYaml.CreateYamlFile(serviceYamlName); err != nil {
		return err
	}

	if len(swaggerFile) > 1 {
		deploymentSwaggerUIYaml := deployFolder + "swagger-ui-deployment.yml"
		ingressSwaggerUIYaml    := deployFolder + "swagger-ui-ingress.yml"
		serviceSwaggerUIYaml    := deployFolder + "swagger-ui-service.yml"
		configMapSwaggerUIYaml  := deployFolder + "swagger-ui-config-map.yml"
		swaggerUIName           := fmt.Sprintf("%s-swagger-ui", functionName)
		swaggerPath             := fmt.Sprintf("/swagger-ui/%s", functionName)

		var swaggerDeploymentYaml d.K8sYaml
		if err := swaggerDeploymentYaml.MountDeployment(swaggerUIName, "swaggerapi/swagger-ui:v3.24.3", swaggerPath); err != nil {
			return err
		}

		if err := swaggerDeploymentYaml.CreateYamlFile(deploymentSwaggerUIYaml); err != nil {
			return err
		}

		var swaggerIngressYaml d.K8sYaml
		if err := swaggerIngressYaml.MountIngress(swaggerUIName, swaggerUIName, swaggerPath, hostnameFlag); err != nil {
			return err
		}

		if err := swaggerIngressYaml.CreateYamlFile(ingressSwaggerUIYaml); err != nil {
			return err
		}

		var swaggerServiceYaml d.K8sYaml
		if err := swaggerServiceYaml.MountService(swaggerUIName, hostnameFlag, false); err != nil {
			return err
		}

		if err := swaggerServiceYaml.CreateYamlFile(serviceSwaggerUIYaml); err != nil {
			return err
		}

		var swaggerConfigMapYaml d.K8sYaml
		if err := swaggerConfigMapYaml.MountConfigMap(swaggerUIName, functionHandler + "/" + swaggerFile); err != nil {
			return err
		}

		if err := swaggerConfigMapYaml.CreateYamlFile(configMapSwaggerUIYaml); err != nil {
			return err
		}
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

func checkSwaggerFileExist(handler string) string {
	swaggerPath := filepath.Join(handler, "swagger.yml")
	if _, err := os.Stat(swaggerPath); err == nil {
		return "swagger.yml"
	}

	swaggerPath = filepath.Join(handler, "swagger.yaml")
	if _, err := os.Stat(swaggerPath); err == nil {
		return "swagger.yaml"
	}

	return ""
}
