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
	exist, err := get.CheckBinary(faasBinaryName, false)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	faasCliPath := config.GetFaasCliPath()
	flagLang, _ := cmd.Flags().GetString("lang")
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	if err := downloadTemplate(faasCliPath, flagLang, verboseFlag); err != nil {
		return err
	}
	
	if err := createFunction(faasCliPath, args[0], flagLang, verboseFlag); err != nil {
		return err
	}
	
	return nil
}

func downloadTemplate(faasCliPath, lang string, verboseFlag bool) error {
	setStore()
	taskDownloadTemplate := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"template", "store", "pull", lang,
		},
		StreamStdio: verboseFlag,
	}

	fmt.Println("Baixando template...")
	res, err := taskDownloadTemplate.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode > 1 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}

func createFunction(faasCliPath, projectName, lang string, verboseFlag bool) error {
	taskCreateFunction := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"new", projectName,
			"--lang", lang,
			"--gateway", "http://gateway.ipaas.localdomain:8080",
			"--prefix", "registry.localdomain:5000",
		},
		StreamStdio: verboseFlag,
	}

	fmt.Println("Criando a " + projectName + "...")
	res, err := taskCreateFunction.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	if err := writeGitignore(); err != nil {
		return err
	}

	if err := addFileEnv(projectName); err != nil {
		return err
	}

	return nil
}

func writeGitignore() error {
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

func addFileEnv(projectName string) error {
	envProjectName := "PROJECT_NAME=" + projectName + "\n"
	envImage := "IMAGE=registry.localdomain:5000/" + projectName + ":latest\n"
	envPort := "PORT=8080\n"
	envDomain := "DOMAIN=ipaas.localdomain\n"

	f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(envProjectName + envImage + envPort + envDomain))
	if err != nil {
		return err
	}

	return nil
}
