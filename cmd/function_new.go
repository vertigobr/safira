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
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	"os"
	"regexp"
)

var newCmd = &cobra.Command{
	Use:     "new FUNCTION_NAME --lang=FUNCTION_LANGUAGE",
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

func PreRunFunctionNew(cmd *cobra.Command, args []string) error {
	flagLang, _ := cmd.Flags().GetString("lang")

	if len(flagLang) == 0 && len(args) < 1 {
		cmd.Help()
		os.Exit(0)
	} else if len(args) < 1 {
		_ = cmd.Help()
		fmt.Println()
		return fmt.Errorf("nome da função não inserido")
	} else if len(flagLang) == 0 {
		return fmt.Errorf("a flag --lang é obrigatória")
	}

	functionName := args[0]
	if err := validateFunctionName(functionName); err != nil {
		return err
	}

	return nil
}

func validateFunctionName(functionName string) error {
	// Regex for RFC-1123 validation:
	// k8s.io/kubernetes/pkg/util/validation/validation.go
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	if matched := validDNS.MatchString(functionName); !matched {
		return fmt.Errorf("o nome da função deve conter apenas caracteres: a-z, 0-9 e traços")
	}

	return nil
}

func runFunctionNew(cmd *cobra.Command, args []string) error {
	fmt.Println(checkDefaultMessage)
	if err := get.CheckFaasCli(); err != nil {
		return err
	}

	faasCliPath := config.GetFaasCliPath()
	flagLang, _ := cmd.Flags().GetString("lang")

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
			"new", projectName,
			"--lang", lang,
			"--gateway", "http://gateway.ipaas.localdomain:8080",
			"--prefix", "registry.localdomain:5000",
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

	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte("deploy\n"))
	if err != nil {
		return err
	}

	return nil
}
