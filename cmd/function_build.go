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

var buildCmd = &cobra.Command{
	Use:   "build -f YAML_FILE",
	Short: "Executa o build de funções",
	Long: `Executa o build de funções`,
	SuggestionsMinimumDistance: 1,
	PreRunE: PreRunFunctionBuild,
	RunE: initFunctionBuild,
}

func init() {
	functionCmd.AddCommand(buildCmd)
}

func initFunctionBuild(cmd *cobra.Command, args []string) error {
	faasCliPath := config.GetFaasCliPath()
	flagYaml, _ := cmd.Flags().GetString("yaml")

	return functionBuild(faasCliPath, flagYaml)
}

func functionBuild(faasCliPath, flagYaml string) error {
	taskFunctionBuild := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"build", "-f", flagYaml,
		},
		StreamStdio: true,
	}

	res, err := taskFunctionBuild.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}

func PreRunFunctionBuild(cmd *cobra.Command, args []string) error {
	flagYaml, _ := cmd.Flags().GetString("yaml")
	if len(flagYaml) == 0 {
		return fmt.Errorf("a flag --yaml/-f é obrigatória")
	}

	return nil
}
