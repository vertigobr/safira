/*
Copyright © Vertigo Tecnologia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	rootCmd.PersistentFlags().Bool("update", false, "Force the deploy to pull a new image")
}

func runFunctionDeploy(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	updateFlag, _ := cmd.Flags().GetBool("update")
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

	if err := functionDeploy(kubectlPath, verboseFlag, updateFlag); err != nil {
		return err
	}

	fmt.Println("\nDeploy realizado com sucesso!")

	return nil
}

func functionDeploy(kubectlPath string, verboseFlag, updateFlag bool) error {
	err := config.SetKubeconfig(clusterName)
	if err != nil {
		return err
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
