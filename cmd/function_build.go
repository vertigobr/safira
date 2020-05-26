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
	Use:     "build -f YAML_FILE",
	Short:   "Executa o build de funções",
	Long:    "Executa o build de funções",
	PreRunE: PreRunFunctionBuild,
	RunE:    runFunctionBuild,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("yaml", "f", "", "Caminho para o yaml de uma função")
}

func PreRunFunctionBuild(cmd *cobra.Command, args []string) error {
	flagYaml, _ := cmd.Flags().GetString("yaml")
	if len(flagYaml) == 0 {
		return fmt.Errorf("a flag --yaml/-f é obrigatória")
	}

	return nil
}

func runFunctionBuild(cmd *cobra.Command, args []string) error {
	fmt.Println(checkDefaultMessage)
	if err := config.CheckFaasCli(); err != nil {
		return err
	}

	faasCliPath := config.GetFaasCliPath()
	flagYaml, _ := cmd.Flags().GetString("yaml")

	if err := functionBuild(faasCliPath, flagYaml); err != nil {
		return err
	}

	if err := functionPush(faasCliPath, flagYaml); err != nil {
		return err
	}

	return nil
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

func functionPush(faasCliPath, flagYaml string) error {
	taskFunctionPush := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"push", "-f", flagYaml,
		},
		StreamStdio: true,
	}

	res, err := taskFunctionPush.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}
