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
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira-libs/pkg/config"
	"github.com/vertigobr/safira-libs/pkg/execute"
	"github.com/vertigobr/safira-libs/pkg/get"
	"os"
	"time"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Levanta uma infraestrutura para ambiente de desenvolvimento.",
	Long: `Levanta uma infraestrutura para ambiente de desenvolvimento com todas as dependências já configuradas.`,
	Run: func(cmd *cobra.Command, args []string) {
		initInfra()
	},
}

func init() {
	infraCmd.AddCommand(upCmd)
}

func checkInfra() {
	if exists, _ := config.ExistsBinary("kubectl"); !exists {
		if err := get.DownloadKubectl(); err != nil {
			panic("Não foi possível baixar o pacote kubectl")
		}
	}

	if exists, _ := config.ExistsBinary("k3d"); !exists {
		if err := get.DownloadK3d(); err != nil {
			panic("Não foi possível baixar o pacote k3d")
		}
	}

	if exists, _ := config.ExistsBinary("helm"); !exists {
		if err := get.DownloadHelm(); err != nil {
			panic("Não foi possível baixar o pacote helm")
		}
	}
}

func createCluster(k3dPath string) error {
	clusterName := "ipaas-local"
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
		StreamStdio: true,
	}

	res, err := taskCreateCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
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
		return fmt.Errorf("exit code %d", resCreateKubeconfig.ExitCode)
	}

	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return err
	}

	return nil
}

func initInfra() {
	checkInfra()
	k3dPath := fmt.Sprintf("%sbin/.%s/%s", config.GetUserDir(), "k3d", "k3d")
	if err := createCluster(k3dPath); err != nil {
		panic(err)
	}
}
