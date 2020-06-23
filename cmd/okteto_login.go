// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login in Okteto",
	Long:  "Login in Okteto",
	RunE:  runOktetoLogin,
	SuggestionsMinimumDistance: 1,
}

func init() {
	oktetoCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("token", "t", "", "API token for authentication")
}

func runOktetoLogin(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	tokenFlag, _ := cmd.Flags().GetString("token")

	err := checkOktetoDependencies(verboseFlag)
	if err != nil {
		return err
	}

	oktetoPath := config.GetOktetoPath()
	if err := loginOkteto(oktetoPath, tokenFlag, verboseFlag); err != nil {
		return err
	}

	return nil
}

func loginOkteto(oktetoPath, tokenFlag string, verboseFlag bool) error {
	args := loginOktetoArgs(tokenFlag)

	taskCreateCluster := execute.Task{
		Command:      oktetoPath,
		Args:         args,
		StreamStdio:  true,
		PrintCommand: verboseFlag,
	}

	fmt.Println("Login Okteto...")
	res, err := taskCreateCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}

func loginOktetoArgs(tokenFlag string) (args []string) {
	args = append(args, "login")

	if len(tokenFlag) > 0 {
		args = append(args, "--token")
		args = append(args, tokenFlag)
	}

	return
}
