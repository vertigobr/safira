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
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira-libs/pkg/config"
	"github.com/vertigobr/safira-libs/pkg/get"
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

func initInfra() {
	checkInfra()
}
