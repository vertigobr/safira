// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/get"
	"gopkg.in/gookit/color.v1"
)

var infraSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Get service credentials",
	Long:  "Get service credentials",
	Example: `To obtain access credentials for some services, run:

    $ safira infra secrets`,
	RunE:                       runInfraSecrets,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraSecretsCmd)
}

func runInfraSecrets(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf(notExistBinary)
	}

	kubectlPath := config.GetKubectlPath()
	if err := getSecrets(kubectlPath, verboseFlag); err != nil {
		return err
	}

	return nil
}

func getSecrets(kubectlPath string, verboseFlag bool) error {
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME")+"/.config/k3d/"+clusterName+"/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("%s It was not possible to export the KUBECONFIG environment variable", color.Red.Text("[!]"))
	}

	taskDeleteCluster := execute.Task{
		Command: kubectlPath,
		Args: []string{
			"--kubeconfig", os.Getenv("KUBECONFIG"),
			"get", "secret", "basic-auth", "-o", `jsonpath={.data.basic-auth-password}`,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("%s Without connection to the cluster", color.Red.Text("[!]"))
	}

	openfaasPassword, _ := base64.StdEncoding.DecodeString(res.Stdout)
	fmt.Println("OpenFaaS - openfaas.ipaas.localdomain:8080")
	fmt.Println("    User: admin")
	fmt.Println("Password: " + string(openfaasPassword))
	fmt.Println("\nKonga - konga.localdomain:8080")
	fmt.Println("    User: admin")
	fmt.Println("Password: admin123")

	return nil
}
