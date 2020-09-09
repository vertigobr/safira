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
	"gopkg.in/gookit/color.v1"
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
	functionDeployCmd.Flags().StringP("namespace", "n", functionsNamespace, fmt.Sprintf("set namespace to deploy (Default: %s)", functionsNamespace))
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
	allFlag, _ := cmd.Flags().GetBool("all-functions")
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
	if allFlag {
		for index := range functions {
			useSha := functions[index].FunctionConfig.Build.UseSha || stack.StackConfig.Build.UseSha
			if err := checkDeployFiles(index, hostnameFlag, namespaceFlag, envFlag, functions[index].Plugins, useSha); err != nil {
				return err
			}

			deployFolder := fmt.Sprintf("deploy/%s/", index)
			if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, index, namespaceFlag, envFlag, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	} else {
		for index, functionArg := range args {
			if checkFunctionExists(args[index], functions) {
				useSha := functions[functionArg].FunctionConfig.Build.UseSha || stack.StackConfig.Build.UseSha
				if err := checkDeployFiles(functionArg, hostnameFlag, namespaceFlag, envFlag, functions[functionArg].Plugins, useSha); err != nil {
					return err
				}

				deployFolder := fmt.Sprintf("deploy/%s/", functionArg)
				if err := deploy(kubectlPath, kubeconfigFlag, deployFolder, functionArg, namespaceFlag, envFlag, verboseFlag, updateFlag); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("nome dá função %s é inválido", functionArg)
			}
		}
	}

	if len(stack.Custom) > 0 {
		for _, path := range stack.Custom {
			if err := deploy(kubectlPath, kubeconfigFlag, path, "", namespaceFlag, envFlag, verboseFlag, updateFlag); err != nil {
				return err
			}
		}
	}

	if swaggerFile := checkSwaggerFileExist(stack.Swagger.File); len(swaggerFile) > 1 {
		if err := deploySwaggerUi(swaggerFile, hostnameFlag, kubectlPath, kubeconfigFlag, envFlag, updateFlag, verboseFlag); err != nil {
			return err
		}
	}

	fmt.Printf("\n%s Deploy successfully completed\n", color.Cyan.Text("[✓]"))

	return nil
}

func deploy(kubectlPath, kubeconfigFlag, deployFolder, functionName, namespace, envFlag string, verboseFlag, updateFlag bool) error {
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
		hasFunction, err := d.CheckFunction(clusterName, functionName, namespace)
		if err != nil {
			return err
		}

		if hasFunction && updateFlag {
			if err := rolloutFunction(kubectlPath, kubeconfig, functionName, namespace, verboseFlag); err != nil {
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
		fmt.Printf("%s Deploying function %s\n", color.Green.Text("[+]"), functionName)
	} else {
		fmt.Printf("%s Deploying custom files %s\n", color.Green.Text("[+]"), deployFolder)
	}

	res, err := taskDeploy.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	if err := addPluginInAnnotations(functionName, namespace, kubeconfig, envFlag, verboseFlag); err != nil {
		return err
	}

	return nil
}

func addPluginInAnnotations(functionName, namespace, kubeconfig, envFlag string, verboseFlag bool) error {
	plugins := ""
	stack, err := s.LoadStackFile(envFlag)
	if err != nil {
		return err
	}

	for pluginName, plugin := range stack.Functions[functionName].Plugins {
		if len(plugin.Type) == 0 || plugin.Type == "service" {
			if len(plugins) > 0 {
				plugins = plugins + ", " + d.GetDeployName(stack, pluginName)
			} else {
				plugins = d.GetDeployName(stack, pluginName)
			}
		}
	}

	if len(stack.Functions[functionName].Plugins) > 0 && len(plugins) > 0 {
		err = d.AddPluginInAnnotationsService(d.GetDeployName(stack, functionName), namespace, plugins, kubeconfig, verboseFlag)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkDeployFiles(functionName, hostnameFlag, namespaceFlag, envFlag string, plugins map[string]s.Plugin, useSha bool) error {
	deployFolder := fmt.Sprintf("./deploy/%s/", functionName)
	functionYamlName := deployFolder + "function.yml"
	ingressYamlName := deployFolder + "ingress.yml"
	//serviceYamlName := deployFolder + "service.yml"
	//functionPath := fmt.Sprintf("/function/%s", functionName)

	if _, err := os.Stat(deployFolder); err != nil {
		if err = os.MkdirAll(deployFolder, 0700); err != nil {
			return err
		}
	}

	var functionYaml d.K8sYaml
	if err := functionYaml.MountFunction(functionName, namespaceFlag, envFlag, useSha); err != nil {
		return err
	}

	if err := functionYaml.CreateYamlFile(functionYamlName); err != nil {
		return err
	}

	var ingressYaml d.K8sYaml
	if err := ingressYaml.MountIngress(functionName, functionName, namespaceFlag, "", hostnameFlag, envFlag); err != nil {
		return err
	}

	if err := ingressYaml.CreateYamlFile(ingressYamlName); err != nil {
		return err
	}

	//var serviceYaml d.K8sYaml
	//if err := serviceYaml.MountService(functionName, hostnameFlag, envFlag, true); err != nil {
	//	return err
	//}
	//
	//if err := serviceYaml.CreateYamlFile(serviceYamlName); err != nil {
	//	return err
	//}

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
			"-n", namespaceFlag,
			"--kubeconfig", kubeconfig,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	if verboseFlag {
		fmt.Printf("%s Resetting the %s function\n", color.Blue.Text("[v]"), functionName)
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

// Receber flag update para executar o rollout no deployment
func deploySwaggerUi(swaggerFile, hostnameFlag, kubectlPath, kubeconfig, envFlag string, update, verboseFlag bool) error {
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
	if err := swaggerDeploymentYaml.MountDeployment(swaggerUIName, "swaggerapi/swagger-ui:v3.24.3", swaggerPath, repoName, envFlag); err != nil {
		return err
	}

	if err := swaggerDeploymentYaml.CreateYamlFile(deploymentSwaggerUIYaml); err != nil {
		return err
	}

	var swaggerIngressYaml d.K8sYaml
	if err := swaggerIngressYaml.MountIngress(swaggerUIName, swaggerUIName, "default", swaggerPath, hostnameFlag, envFlag); err != nil {
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
	if err := swaggerConfigMapYaml.MountConfigMap(swaggerUIName, swaggerFile, repoName, envFlag); err != nil {
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

	fmt.Printf("%s Deploying Swagger UI\n", color.Green.Text("[+]"))

	res, err := taskDeploySwaggerUi.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	if update {
		taskRolloutSwagger := execute.Task{
			Command: kubectlPath,
			Args: []string{
				"rollout", "restart", "deployments", swaggerUIName,
				"--kubeconfig", kubeconfig,
			},
			StreamStdio:  verboseFlag,
			PrintCommand: verboseFlag,
		}

		res, err := taskRolloutSwagger.Execute()
		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}
	}

	return nil
}

func checkSwaggerFileExist(fileName string) string {
	if len(fileName) > 0 {
		swaggerPath := filepath.Join(fileName)
		if _, err := os.Stat(swaggerPath); err == nil {
			return fileName
		}
	} else {
		swaggerPath := filepath.Join("swagger.yml")
		if _, err := os.Stat(swaggerPath); err == nil {
			return "swagger.yml"
		}

		swaggerPath = filepath.Join("swagger.yaml")
		if _, err := os.Stat(swaggerPath); err == nil {
			return "swagger.yaml"
		}
	}

	return ""
}
