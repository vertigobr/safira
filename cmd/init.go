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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
	"os"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Sincroniza todas as depências",
	Long:  "Sincroniza todas as depências para uso do Safira",
	PreRunE: PreRunInit,
	RunE: runInit,
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func PreRunInit(cmd *cobra.Command, args []string) error {
	home := os.Getenv("HOME")

	if home == "/root" {
		return fmt.Errorf("comando init executado de forma inválida, execute: \n\n\t" + safiraInit)
	}

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	fmt.Println("Verificando dependências...")
	if _, err := checkAllBinaries(verboseFlag); err != nil {
		return err
	}

	if err := checkHosts(verboseFlag); err != nil {
		return err
	}

	fmt.Println("\nTodas as dependências resolvidas!")

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

func checkHosts(verboseFlag bool) error {
	host := "127.0.0.1 registry.localdomain ipaas.localdomain konga.localdomain gateway.ipaas.localdomain"
	hostsFile := "/etc/hosts"
	if verboseFlag {
		fmt.Println("[+] Verificando hosts")
	}

	fileRead, err := os.OpenFile(hostsFile, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	scanner := bufio.NewScanner(fileRead)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), host) {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fileWrite, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	_, err = fileWrite.Write([]byte(host + "\n"))
	if err != nil {
		return err
	}

	if verboseFlag {
		fmt.Println("[+] Gravado com sucesso no /etc/hosts")
	}

	return nil
}
