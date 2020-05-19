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
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Derruba uma infraestrutura provisionada anteriormente.",
	Long: `Derruba uma infraestrutura provisionada anteriormente para desenvolvimento.`,
	SuggestionsMinimumDistance: 1,
	Run: func(cmd *cobra.Command, args []string) {
		downInfra()
	},
}

func init() {
	infraCmd.AddCommand(downCmd)
}

func deleteCluster(k3dPath string) error {
	taskDeleteCluster := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"delete",
			"-n", clusterName,
		},
		StreamStdio: false,
	}

	fmt.Println("Destruindo cluster local...")
	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}

func downInfra() {
	checkInfra()
	k3dPath := config.GetK3dPath()

	if err := deleteCluster(k3dPath); err != nil {
		fmt.Println("\nNenhum cluster encontrado")
		return
	}

	fmt.Println("\nCluster destruído com sucesso!")
}
