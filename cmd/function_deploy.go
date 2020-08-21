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
	"github.com/vertigobr/safira/pkg/utils"
)

var functionDeployCmd = &cobra.Command{
	Use:   "deploy [FUNCTION_NAME]",
	Short: "Deploy functions",
	Long:  "Deploy functions",
	Example: `If you want to deploy a function's Docker image, run:

    $ safira function deploy function-name

or if you want to deploy the Docker image of all the functions, execute:

    $ safira function deploy -A`,
	PreRunE:                    preRunFunctionDeploy,
	RunE:                       runFunctionDeploy,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionDeployCmd)
	functionDeployCmd.Flags().Bool("update", false, "force the deploy to pull a new image (Default: false)")
	functionDeployCmd.Flags().BoolP("all-functions", "A", false, "deploy all functions")
	functionDeployCmd.Flags().String("kubeconfig", "", "set kubeconfig to deploy")
	functionDeployCmd.Flags().String("hostname", "", "set hostname to deploy")
	functionDeployCmd.Flags().StringP("namespace", "n", "", fmt.Sprintf("set namespace to deploy (Default: %s)", functionsNamespace))
	functionDeployCmd.Flags().StringP("env", "e", "", "Set stack env file")
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
	all, _ := cmd.Flags().GetBool("all-functions")
	kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")
	hostnameFlag, _ := cmd.Flags().GetString("hostname")
	namespaceFlag, _ := cmd.Flags().GetString("namespace")
	envFlag, _ := cmd.Flags().GetString("env")

	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()
	stack, err := s.LoadStackFile(envFlag)
	if err != nil {
		return err
	}

	functions := stack.Functions
	if all {
		for index := range functions {
			if err := checkDeployFiles(index, hostnameFlag, namespaceFlag, envFlag, functions[index].Plugins); err != nil {
				return err
			}

			deployFolder := fmt.Sprintf("deploy/%s/", index)
			if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, index, namespaceFlag, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {

				if err := checkDeployFiles(functionArg, hostnameFlag, namespaceFlag, envFlag, functions[functionArg].Plugins); err != nil {
					return err
				}

				deployFolder := fmt.Sprintf("deploy/%s/", functionArg)
				if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, functionArg, namespaceFlag, verboseFlag, updateFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	if len(stack.Custom) > 0 {
		for _, path := range stack.Custom {
			if err := deploy(kubectlPath, kubeconfigFlag, path, "", namespaceFlag, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	}

	if swaggerFile := checkSwaggerFileExist(); len(swaggerFile) > 1 {
		if err := deploySwaggerUi(swaggerFile, hostnameFlag, kubectlPath, kubeconfigFlag, envFlag, verboseFlag); err != nil {
			return err
		}
	}

	fmt.Println("\nDeploy realizado com sucesso!")

	return nil
}

func deploy(kubectlPath, kubeconfigFlag, deployFolder, functionName, namespaceFlag string, verboseFlag, updateFlag bool) error {
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
		hasFunction, err := d.CheckFunction(clusterName, functionName, getNamespaceDeploy(namespaceFlag))
		if err != nil {
			return err
		}

		if hasFunction && updateFlag {
			if err := rolloutFunction(kubectlPath, kubeconfig, functionName, namespaceFlag, verboseFlag); err != nil {
				return err
			}
		}
	}

	taskDeploy := execute.Task{
		Command: kubectlPath,
		Args: []string{
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

func checkDeployFiles(functionName, hostnameFlag, namespaceFlag, envFlag string, plugins map[string]s.Plugin) error {
	deployFolder := fmt.Sprintf("./deploy/%s/", functionName)
	functionYamlName := deployFolder + "function.yml"
	ingressYamlName := deployFolder + "ingress.yml"
	serviceYamlName := deployFolder + "service.yml"
	functionPath := fmt.Sprintf("/function/%s", functionName)

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	var functionYaml d.K8sYaml
	if err := functionYaml.MountFunction(functionName, getNamespaceDeploy(namespaceFlag), envFlag); err != nil {
		return err
	}

	if err := functionYaml.CreateYamlFile(functionYamlName); err != nil {
		return err
	}

	var ingressYaml d.K8sYaml
	if err := ingressYaml.MountIngress(functionName, functionName, functionPath, hostnameFlag, envFlag); err != nil {
		return err
	}

	if err := ingressYaml.CreateYamlFile(ingressYamlName); err != nil {
		return err
	}

	var serviceYaml d.K8sYaml
	if err := serviceYaml.MountService(functionName, hostnameFlag, envFlag, true); err != nil {
		return err
	}

	if err := serviceYaml.CreateYamlFile(serviceYamlName); err != nil {
		return err
	}

	for pluginName := range plugins {
		pluginYamlName := deployFolder + pluginName + ".yml"
		var pluginYaml d.K8sYaml
		if err := pluginYaml.MountKongPlugin(functionName, pluginName, functionsNamespace, envFlag); err != nil {
			return err
		}

		if err := pluginYaml.CreateYamlFile(pluginYamlName); err != nil {
			return err
		}
	}

	return nil
}

func rolloutFunction(kubectlPath, kubeconfig, functionName, namespaceFlag string, verboseFlag bool) error {
	taskRemoveFunction := execute.Task{
		Command: kubectlPath,
		Args: []string{
			"rollout", "restart", "deployments", functionName,
			"-n", getNamespaceDeploy(namespaceFlag),
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

func deploySwaggerUi(swaggerFile, hostnameFlag, kubectlPath, kubeconfig, envFlag string, verboseFlag bool) error {
	deployFolder := "./deploy/swagger-ui/"
	repoName, err := utils.GetCurrentFolder()
	if err != nil {
		return err
	}

	deploymentSwaggerUIYaml := deployFolder + "swagger-ui-deployment.yml"
	ingressSwaggerUIYaml := deployFolder + "swagger-ui-ingress.yml"
	serviceSwaggerUIYaml := deployFolder + "swagger-ui-service.yml"
	configMapSwaggerUIYaml := deployFolder + "swagger-ui-config-map.yml"
	swaggerUIName := fmt.Sprintf("%s-swagger-ui", repoName)
	swaggerPath := fmt.Sprintf("/swagger-ui/%s", repoName)

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	var swaggerDeploymentYaml d.K8sYaml
	if err := swaggerDeploymentYaml.MountDeployment(swaggerUIName, "swaggerapi/swagger-ui:v3.24.3", swaggerPath); err != nil {
		return err
	}

	if err := swaggerDeploymentYaml.CreateYamlFile(deploymentSwaggerUIYaml); err != nil {
		return err
	}

	var swaggerIngressYaml d.K8sYaml
	if err := swaggerIngressYaml.MountIngress(swaggerUIName, swaggerUIName, swaggerPath, hostnameFlag, envFlag); err != nil {
		return err
	}

	if err := swaggerIngressYaml.CreateYamlFile(ingressSwaggerUIYaml); err != nil {
		return err
	}

	var swaggerServiceYaml d.K8sYaml
	if err := swaggerServiceYaml.MountService(swaggerUIName, hostnameFlag, envFlag, false); err != nil {
		return err
	}

	if err := swaggerServiceYaml.CreateYamlFile(serviceSwaggerUIYaml); err != nil {
		return err
	}

	var swaggerConfigMapYaml d.K8sYaml
	if err := swaggerConfigMapYaml.MountConfigMap(swaggerUIName, swaggerFile); err != nil {
		return err
	}

	if err := swaggerConfigMapYaml.CreateYamlFile(configMapSwaggerUIYaml); err != nil {
		return err
	}

	taskDeploySwaggerUi := execute.Task{
		Command: kubectlPath,
		Args: []string{
			"apply", "--wait",
			"--kubeconfig", kubeconfig,
			"-f", deployFolder,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Executando deploy do Swagger UI...")

	res, err := taskDeploySwaggerUi.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}

func checkSwaggerFileExist() string {
	swaggerPath := filepath.Join("swagger.yml")
	if _, err := os.Stat(swaggerPath); err == nil {
		return "swagger.yml"
	}

	swaggerPath = filepath.Join("swagger.yaml")
	if _, err := os.Stat(swaggerPath); err == nil {
		return "swagger.yaml"
	}

	return ""
}

func getNamespaceDeploy(namespaceFlag string) string {
	if len(namespaceFlag) > 0 {
		return namespaceFlag
	}

	return functionsNamespace
}
