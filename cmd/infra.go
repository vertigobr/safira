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
)

const clusterName = "vertigo-ipaas"

var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Responsável por gerenciar a infraestrutura",
	Long:  "Responsável por gerenciar a infraestrutura em ambiente local", //  ou em nuvem
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(infraCmd)

	infraCmd.PersistentFlags().String(
		"env",
		"local",
		"Recebe o ambiente aonde será provisionado o cluster Kubernetes.",
	)
}

func checkInfra() error {
	fmt.Println(checkDefaultMessage)

	if err := config.CheckKubectl(); err != nil {
		return err
	}

	if err := config.CheckK3d(); err != nil {
		return err
	}

	if err := config.CheckHelm(); err != nil {
		return err
	}

	return nil
}
