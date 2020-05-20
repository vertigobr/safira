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
	"errors"
	"fmt"
	"github.com/vertigobr/safira-libs/pkg/execute"
	"github.com/vertigobr/safira-libs/pkg/repository"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [REPOSITORY] [NAME]",
	Args: validArgs,
	Short: "Faz download de templates",
	Long: `Faz download de templates podendo declarar templates oficiais ou privados`,
	Example: `
  safira template pull java11 project-name
  safira template pull https://github.com/owner/repository project-name
`,
	SuggestionsMinimumDistance: 1,
	RunE: func(cmd *cobra.Command, args []string) error {
		return templatePull(args)
	},
}

func init() {
	templateCmd.AddCommand(pullCmd)
}

func validArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("repositório não inserido")
	}

	if len(args) < 2 {
		return errors.New("nome do projeto não inserido")
	}

	return nil
}

func templatePull(args []string) error {
	repo, err := repository.GetRepositoryURL(args[0])
	if err != nil {
		return err
	}

	taskPull := execute.Task{
		Command:     "git",
		Args:        []string{
			"clone", repo, args[1],
		},
		StreamStdio: false,
	}

	fmt.Println("Baixando template...")
	res, err := taskPull.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}
