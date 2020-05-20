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
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira-libs/pkg/config"
	"github.com/vertigobr/safira-libs/pkg/execute"
	"net/url"
)

var pullCmd = &cobra.Command{
	Use:   "pull [REPOSITORY]",
	Args: validArgs,
	Short: "Faz download de templates",
	Long: `Faz download de templates, podendo declarar templates oficiais (utilizando apenas o nome) ou privados (passando a URL)`,
	Example: `
  safira template pull vtg-ipaas-java11
  safira template pull https://github.com/owner/repository
`,
	SuggestionsMinimumDistance: 1,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initTemplatePull(args[0])
	},
}

func validArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("repositório não inserido")
	}

	return nil
}

func init() {
	templateCmd.AddCommand(pullCmd)
}

func initTemplatePull(repository string) error {
	faasCliPath := config.GetFaasCliPath()

	repo, official := checkRepositoryTemplate(repository)

	if err := templatePull(faasCliPath, repo, official); err != nil {
		return err
	}

	return nil
}

func templatePull(faasCliPath, repo string, official bool) error {
	checkOpenFaas()
	var taskPull execute.Task

	if official {
		taskPull = execute.Task{
			Command:     faasCliPath,
			Args:        []string{
				"template", "store", "pull", repo,
			},
			StreamStdio: true,
		}
	} else {
		taskPull = execute.Task{
			Command:     faasCliPath,
			Args:        []string{
				"template", "pull", repo,
			},
			StreamStdio: true,
		}
	}

	res, err := taskPull.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode > 1 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}

func checkRepositoryTemplate(repository string) (string, bool) {
	u, _ := url.Parse(repository)

	if u.Host == "" {
		return repository, true
	} else {
		return repository, false
	}
}
