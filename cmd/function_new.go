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

var newCmd = &cobra.Command{
	Use:     "new FUNCTION_NAME --lang=FUNCTION_LANGUAGE",
	Args:    validArgsFunctionNew,
	Short:   "Cria uma nova função na pasta atual",
	Long:    "Cria uma nova função hello-world baseada na linguagem inserida",
	Example: "safira function new project-name --lang=java",
	PreRunE: PreRunFunctionNew,
	RunE:    runFunctionNew,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(newCmd)

	newCmd.Flags().String("lang", "", "Linguagem para criação do template")
}

func validArgsFunctionNew(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		_ = cmd.Help()
		fmt.Println()
		return fmt.Errorf("nome da função não inserido")
	}

	return nil
}

func PreRunFunctionNew(cmd *cobra.Command, args []string) error {
	flagLang, _ := cmd.Flags().GetString("lang")
	if len(flagLang) == 0 {
		return fmt.Errorf("a flag --lang é obrigatória")
	}

	return nil
}

func runFunctionNew(cmd *cobra.Command, args []string) error {
	faasCliPath := config.GetFaasCliPath()
	flagLang, _ := cmd.Flags().GetString("lang")
	checkOpenFaas()

	if err := downloadTemplate(faasCliPath, flagLang); err != nil {
		return err
	}
	
	if err := createFunction(faasCliPath, args[0], flagLang); err != nil {
		return err
	}
	
	return nil
}

func downloadTemplate(faasCliPath, lang string) error {
	setStore()
	taskDownloadTemplate := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"template", "store", "pull", lang,
		},
		StreamStdio: true,
	}

	fmt.Println("Baixando template...")
	res, err := taskDownloadTemplate.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode > 1 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}

func createFunction(faasCliPath, projectName, lang string) error {
	taskCreateFunction := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"new", projectName, "--lang", lang,
		},
		StreamStdio: true,
	}

	fmt.Println("Criando template...")
	res, err := taskCreateFunction.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}
