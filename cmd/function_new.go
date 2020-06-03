// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	"github.com/vertigobr/safira/pkg/stack"
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
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	exist, err := get.CheckBinary(faasBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	faasCliPath := config.GetFaasCliPath()
	flagLang, _ := cmd.Flags().GetString("lang")
	functionName := args[0]

	if err := downloadTemplate(faasCliPath, flagLang, verboseFlag); err != nil {
		return err
	}

	function := stack.Function{
		Name:    functionName,
		Lang:    flagLang,
		Handler: "./" + functionName,
		Image:   "registry.localdomain:5000/" + functionName + ":latest",
	}

	if _, err = os.Stat("stack.yml"); err != nil {
		if err := stack.CreateTemplate(function); err != nil {
			return err
		}
	} else {
		if err := stack.AppendFunction(function); err != nil {
			return err
		}
	}

	if err := createFunction(faasCliPath, functionName, flagLang, verboseFlag); err != nil {
		return err
	}

	fmt.Println("\nFunction " + functionName + " criada com sucesso!")

	return nil
}

func downloadTemplate(faasCliPath, lang string, verboseFlag bool) error {
	taskDownloadTemplate := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"template", "store", "pull", lang, "--url", faasTemplateStoreURL,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
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

func createFunction(faasCliPath, functionName, lang string, verboseFlag bool) error {
	taskCreateFunction := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"new", functionName,
			"--lang", lang,
			"--gateway", "http://gateway.ipaas.localdomain:8080",
			"--prefix", "registry.localdomain:5000",
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Criando function " + functionName + "...")
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

	if err := deleteYamlFunction(functionName); err != nil {
		return fmt.Errorf("error ao remover yaml da função gerada %s.yml: %s", functionName, err.Error())
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

func deleteYamlFunction(functionName string) error {
	return os.Remove(functionName + ".yml")
}
