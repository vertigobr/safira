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
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	"os"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Obtém informações de acesso as plataformas",
	Long:  "Obtém informações de acesso as plataformas do Vertigo iPaaS",
	RunE:  runInfraSecrets,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(secretsCmd)
}

func runInfraSecrets(cmd *cobra.Command, args []string) error {
	fmt.Println(checkDefaultMessage)
	if err := get.CheckKubectl(); err != nil {
		return err
	}

	kubectlPath := config.GetKubectlPath()
	if err := getSecrets(kubectlPath); err != nil {
		return err
	}

	return nil
}

func getSecrets(kubectlPath string) error {
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível adicionar a variável de ambiente KUBECONFIG")
	}

	taskDeleteCluster := execute.Task{
		Command:     kubectlPath,
		Args:        []string{
			"--kubeconfig", os.Getenv("KUBECONFIG"),
			"get", "secret", "basic-auth", "-o",`jsonpath={.data.basic-auth-password}`,
		},
		StreamStdio: false,
	}

	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		fmt.Println()
		return fmt.Errorf("sem conexão com o cluster")
	}

	openfaasPassword, _ := base64.StdEncoding.DecodeString(res.Stdout)
	fmt.Println("\nOpenFaaS - gateway.ipaas.localdomain:8080")
	fmt.Println("    User: admin")
	fmt.Println("Password: " + string(openfaasPassword))

	return nil
}
