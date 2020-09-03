// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
	"gopkg.in/gookit/color.v1"
)

var initCmd = &cobra.Command{
	Use:                        "init",
	Short:                      "Synchronizes all the dependencies necessary for the sapphire to function",
	Long:                       "Synchronizes all the dependencies necessary for the sapphire to function",
	PreRunE:                    PreRunInit,
	RunE:                       runInit,
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func PreRunInit(_ *cobra.Command, _ []string) error {
	if os.Getuid() != 0 {
		return fmt.Errorf("%s Init command executed invalidly, execute: \n\n\t%s", color.Red.Text("[!]"), safiraInit)
	}

	return nil
}

func runInit(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	fmt.Printf("%s Checking dependencies\n", color.Green.Text("[+]"))
	if _, err := checkAllBinaries(verboseFlag); err != nil {
		return err
	}

	if err := checkHosts(verboseFlag); err != nil {
		return err
	}

	fmt.Printf("\n%s All dependencies resolved\n", color.Cyan.Text("[✓]"))

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

	_, err = get.CheckBinary(oktetoBinaryName, true, verboseFlag)
	if err != nil {
		return true, err
	}

	return true, nil
}

func checkHosts(verboseFlag bool) error {
	host := "127.0.0.1 registry.localdomain ipaas.localdomain konga.localdomain openfaas.ipaas.localdomain editor.localdomain"
	hostsFile := "/etc/hosts"

	if verboseFlag {
		fmt.Printf("%s Checking hosts file\n", color.Blue.Text("[v]"))
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
		fmt.Printf("%s Successfully saved to hosts file\n", color.Blue.Text("[v]"))
	}

	return nil
}
