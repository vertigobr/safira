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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista os templates oficiais do Vertigo iPaaS",
	Long:  "Lista os templates oficiais do Vertigo iPaaS",
	RunE:  runTemplateList,
	SuggestionsMinimumDistance: 1,
}

func init() {
	templateCmd.AddCommand(listCmd)
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	faasCliPath := config.GetFaasCliPath()
	checkOpenFaas()

	if err := templateList(faasCliPath); err != nil {
		return err
	}

	return nil
}

func templateList(faasCliPath string) error {
	setStore()

	taskList := execute.Task{
		Command:     faasCliPath,
		Args:        []string{
			"template", "store", "list",
		},
		StreamStdio: true,
	}

	res, err := taskList.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	return nil
}
