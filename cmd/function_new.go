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
	"gopkg.in/gookit/color.v1"
)

var functionNewCmd = &cobra.Command{
	Use:   "new [FUNCTION_NAME] --lang=[FUNCTION_LANGUAGE or TEMPLATE_NAME]",
	Short: "Creates a new function",
	Long:  "Creates a new hello-world function based on the inserted language",
	Example: `to create a new function, run:

    $ safira function new project-name --lang=template-name`,
	PreRunE:                    preRunFunctionNew,
	RunE:                       runFunctionNew,
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
		_ = cmd.Help()
		os.Exit(0)
	} else if len(args) < 1 {
		_ = cmd.Help()
		fmt.Println()
		return fmt.Errorf("%s Function name is required", color.Red.Text("[!]"))
	} else if len(flagLang) == 0 {
		return fmt.Errorf("%s The %s flag is required", color.Red.Text("[!]"), color.Bold.Text("lang"))
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
		return fmt.Errorf("%s Function name must contain only characters: a-z, 0-9 and dashes", color.Red.Text("[!]"))
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

	if _, err = os.Stat(stack.GetStackFileName()); err != nil {
		if err := stack.CreateTemplate(function); err != nil {
			return err
		}
	} else {
		if err := stack.AppendFunction(function); err != nil {
			return err
		}
	}

	if _, err = os.Stat(ci.GitlabCiFileName); err != nil {
		if err := ci.CreateFile(); err != nil {
			return err
		}
	} // else {
	//	if err := ci.AppendFunction(functionName); err != nil {
	//		return err
	//	}
	//}

	if err := createFunction(faasCliPath, functionName, flagLang, verboseFlag); err != nil {
		return err
	}

	fmt.Printf("%s Function %s successfully created\n", color.Green.Text("[+]"), functionName)

	return nil
}

func createFunction(faasCliPath, functionName, lang string, verboseFlag bool) error {
	taskCreateFunction := execute.Task{
		Command: faasCliPath,
		Args: []string{
			"new", functionName,
			"--lang", lang,
			"--gateway", "http://gateway.ipaas.localdomain:8080",
			"--prefix", "registry.localdomain:5000",
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	fmt.Printf("%s Creating function %s\n", color.Green.Text("[+]"), functionName)
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
		return fmt.Errorf("%s Error when removing function %s", color.Green.Text("[+]"), functionName)
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
