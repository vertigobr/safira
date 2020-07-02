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
)

var infraSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Obtém informações de acesso as plataformas",
	Long:  "Obtém informações de acesso as plataformas do Vertigo iPaaS",
	RunE:  runInfraSecrets,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraSecretsCmd)
}

func runInfraSecrets(cmd *cobra.Command, args []string) error {
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
	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível adicionar a variável de ambiente KUBECONFIG")
	}

	taskDeleteCluster := execute.Task{
		Command:     kubectlPath,
		Args:        []string{
			"--kubeconfig", os.Getenv("KUBECONFIG"),
			"get", "secret", "basic-auth", "-o",`jsonpath={.data.basic-auth-password}`,
		},
		StreamStdio:  verboseFlag,
		PrintCommand: verboseFlag,
	}

	res, err := taskDeleteCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		fmt.Println()
		return fmt.Errorf("sem conexão com o cluster")
	}

	openfaasPassword, _ := base64.StdEncoding.DecodeString(res.Stdout)
	fmt.Println("OpenFaaS - gateway.ipaas.localdomain:8080")
	fmt.Println("    User: admin")
	fmt.Println("Password: " + string(openfaasPassword))
	fmt.Println("\nKonga - konga.localdomain:8080")
	fmt.Println("    User: admin")
	fmt.Println("Password: admin123")

	return nil
}
