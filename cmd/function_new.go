// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/ci"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	"github.com/vertigobr/safira/pkg/stack"
)

var functionNewCmd = &cobra.Command{
	Use:     "new [FUNCTION_NAME] --lang=FUNCTION_LANGUAGE",
	Short:   "Creates a new function",
	Long:    "Creates a new hello-world function based on the inserted language",
	Example: "safira function new project-name --lang=java",
	PreRunE: preRunFunctionNew,
	RunE:    runFunctionNew,
	SuggestionsMinimumDistance: 1,
}

func init() {
	functionCmd.AddCommand(functionNewCmd)

	functionNewCmd.Flags().String("lang", "", "template name")
	functionNewCmd.Flags().Bool("update-template", false, "update template folder")
}

func preRunFunctionNew(cmd *cobra.Command, args []string) error {
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
	updateTemplateFlag, _ := cmd.Flags().GetBool("update-template")
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

	if err := get.DownloadTemplate(faasTemplateRepo, updateTemplateFlag, verboseFlag); err != nil {
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

	if _, err = os.Stat(".gitlab-ci.yml"); err != nil {
		if err := ci.CreateFile(); err != nil {
			return err
		}
	}

	if err := createFunction(faasCliPath, functionName, flagLang, verboseFlag); err != nil {
		return err
	}

	fmt.Println("\nFunction " + functionName + " criada com sucesso!")

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
	gitIgnoreFile := ".gitignore"
	fileRead, err := os.Open(gitIgnoreFile)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	scanner := bufio.NewScanner(fileRead)
	for scanner.Scan() {
		if scanner.Text() == "deploy" {
			return nil
		}
	}

	fileWrite, err := os.OpenFile(gitIgnoreFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	_, err = fileWrite.Write([]byte("deploy\n"))
	if err != nil {
		return err
	}

	return nil
}

func deleteYamlFunction(functionName string) error {
	return os.Remove(functionName + ".yml")
}
