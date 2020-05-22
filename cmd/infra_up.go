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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira-libs/pkg/config"
	"github.com/vertigobr/safira-libs/pkg/execute"
	"os"
	"time"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Levanta uma infraestrutura para ambiente de desenvolvimento",
	Long: "Levanta uma infraestrutura para ambiente de desenvolvimento com todas as dependências já configuradas",
	SuggestionsMinimumDistance: 1,
	RunE: initInfraUp,
}

func init() {
	infraCmd.AddCommand(upCmd)
}

func initInfraUp(cmd *cobra.Command, args []string) error {
	if err := checkInfra(); err != nil {
		return err
	}

	k3dPath := config.GetK3dPath()
	helmPath := config.GetHelmPath()

	if err := createCluster(k3dPath); err != nil {
		return err
	}

	if err := helmUpgrade(helmPath); err != nil {
		return err
	}

	fmt.Println("\nCluster criado com sucesso!")
	fmt.Println("Konga    - konga.localdomain:8080")
	fmt.Println("Gateway  - ipaas.localdomain:8080")
	fmt.Println("OpenFaaS - gateway.ipaas.localdomain:8080")

	return nil
}

func createCluster(k3dPath string) error {
	taskCreateCluster := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"create",
			"-n", clusterName,
			"--enable-registry",
			"--registry-name", "registry.localdomain",
			"--publish", "8080:32080",
			"-server-arg", "--no-deploy=traefik",
			"-server-arg", "--no-deploy=servicelb",
		},
		StreamStdio: false,
	}

	fmt.Println("Provisionando cluster local...")
	res, err := taskCreateCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return errors.New("\nCluster local já está levantado!")
	}

	time.Sleep(time.Second * 10)

	taskCreateKubeconfig := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"get-kubeconfig",
			"-n", clusterName,
		},
		StreamStdio: false,
	}

	resCreateKubeconfig, err := taskCreateKubeconfig.Execute()
	if err != nil {
		return err
	}

	if resCreateKubeconfig.ExitCode != 0 {
		return errors.New("\nFalha na exportação do KUBECONFIG, tente novamente!")
	}

	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível adicionar a variável de ambiente KUBECONFIG")
	}

	return nil
}

func helmUpgrade(helmPath string) error {
	taskRepoAdd := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"repo", "add", "vtg-ipaas",
			"https://vertigobr.gitlab.io/ipaas/vtg-ipaas-chart",
		},
		StreamStdio: false,
	}

	fmt.Println("Instalando o Vertigo iPaaS...")
	resRepoAdd, err := taskRepoAdd.Execute()
	if err != nil {
		return err
	}

	if resRepoAdd.ExitCode != 0 {
		return fmt.Errorf("exit code %d", resRepoAdd.ExitCode)
	}

	taskRepoUpdate := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"repo", "update",
		},
		StreamStdio: false,
	}

	resRepoUpdate, err := taskRepoUpdate.Execute()
	if err != nil {
		return err
	}

	if resRepoUpdate.ExitCode != 0 {
		return fmt.Errorf("exit code %d", resRepoUpdate.ExitCode)
	}

	taskUpgrade := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"upgrade", "-i",
			"--kubeconfig", os.Getenv("KUBECONFIG"),
			"-f", "https://gist.githubusercontent.com/kyfelipe/1db230e45d14213ea5ca375aa74057e4/raw/859a26cbb488012f0af6520b8dab253abf2fd97e/k3d.yaml",
			"vtg-ipaas", "vtg-ipaas/vtg-ipaas",
		},
		StreamStdio: false,
	}

	resUpgrade, err := taskUpgrade.Execute()
	if err != nil {
		return err
	}

	if resUpgrade.ExitCode != 0 {
		return fmt.Errorf("exit code %d", resUpgrade.ExitCode)
	}

	return nil
}
