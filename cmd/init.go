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
	"github.com/vertigobr/safira/pkg/get"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Sincroniza todas as depências",
	Long:  "Sincroniza todas as depências para uso do Safira",
	RunE: runInit,
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	_, err := checkAllBinaries(verboseFlag)
	if err != nil {
		return err
	}

	fmt.Println("Todas as dependências resolvidas!")

	return nil
}

func checkAllBinaries(verboseFlag bool) (bool, error) {
	_, err := get.CheckBinary(kubectlBinaryName, true, verboseFlag)
	if err != nil {
		return true, err
	}

	_, err = get.CheckBinary(k3dBinaryName, true, verboseFlag)
	if err != nil {
		return true, err
	}

	_, err = get.CheckBinary(helmBinaryName, true, verboseFlag)
	if err != nil {
		return true, err
	}

	_, err = get.CheckBinary(faasBinaryName, true, verboseFlag)
	if err != nil {
		return true, err
	}

	return true, nil
}

