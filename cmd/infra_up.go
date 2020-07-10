// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
)

var infraUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Provides a local cluster for the development environment",
	Long:  "Provides a local cluster for the development environment",
	Example: `To provision the local cluster for development, run:

    $ safira infra up`,
	RunE:  runInfraUp,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraUpCmd)
}

func runInfraUp(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	err := checkInfra(verboseFlag)
	if err != nil {
		return err
	}

	k3dPath := config.GetK3dPath()
	helmPath := config.GetHelmPath()
	kubectlPath := config.GetKubectlPath()

	if err := createCluster(k3dPath, verboseFlag); err != nil {
		return err
	}

	if err := createNamespace(kubectlPath, verboseFlag); err != nil {
		return err
	}

	if err := helmUpgrade(helmPath, verboseFlag); err != nil {
		return err
	}

	fmt.Println()
	fmt.Print(`Cluster criado com sucesso!
Konga    - konga.localdomain:8080
Gateway  - ipaas.localdomain:8080
OpenFaaS - gateway.ipaas.localdomain:8080
Editor   - editor.localdomain:8080

Para acesso ao cluster:
export KUBECONFIG=$(safira infra get-kubeconfig)
`)

	return nil
}

func createCluster(k3dPath string, verboseFlag bool) error {
	taskCreateCluster := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"create",
			"-n", clusterName,
			"--enable-registry",
			"--registry-name", "registry.localdomain",
			"--publish", "8080:32080",
			"-server-arg", "--no-deploy=traefik",
			"-server-arg", "--no-deploy=servicelb",
		},
		StreamStdio: verboseFlag,
	}

	fmt.Println("Provisionando cluster local...")
	res, err := taskCreateCluster.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	time.Sleep(time.Second * 10)

	taskCreateKubeconfig := execute.Task{
		Command:     k3dPath,
		Args:        []string{
			"get-kubeconfig",
			"-n", clusterName,
		},
		StreamStdio: false,
	}

	resCreateKubeconfig, err := taskCreateKubeconfig.Execute()
	if err != nil {
		return err
	}

	if resCreateKubeconfig.ExitCode != 0 {
		return errors.New("\nFalha na exportação do KUBECONFIG, tente novamente!")
	}

	if err := os.Setenv("KUBECONFIG", os.Getenv("HOME") + "/.config/k3d/" + clusterName + "/kubeconfig.yaml"); err != nil {
		return fmt.Errorf("não foi possível adicionar a variável de ambiente KUBECONFIG")
	}

	return nil
}

func createNamespace(kubectl string, verboseFlag bool) error {
	taskCreateNamespace := execute.Task{
		Command:     kubectl,
		Args:        []string{
			"create", "namespace", functionsNamespace,
		},
		StreamStdio: verboseFlag,
	}

	if verboseFlag {
		fmt.Println("[+] Criando namespace ipaas-fn")
	}

	resCreateNamespace, err := taskCreateNamespace.Execute()
	if err != nil {
		return err
	}

	if resCreateNamespace.ExitCode != 0 {
		return fmt.Errorf(resCreateNamespace.Stderr)
	}

	return nil
}

func helmUpgrade(helmPath string, verboseFlag bool) error {
	taskRepoAdd := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"repo", "add", "vtg-ipaas",
			"https://vertigobr.gitlab.io/ipaas/vtg-ipaas-chart",
		},
		StreamStdio: verboseFlag,
	}

	fmt.Println("Instalando o Vertigo iPaaS...")
	resRepoAdd, err := taskRepoAdd.Execute()
	if err != nil {
		return err
	}

	if resRepoAdd.ExitCode != 0 {
		return fmt.Errorf(resRepoAdd.Stderr)
	}

	taskRepoUpdate := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"repo", "update",
		},
		StreamStdio: verboseFlag,
	}

	resRepoUpdate, err := taskRepoUpdate.Execute()
	if err != nil {
		return err
	}

	if resRepoUpdate.ExitCode != 0 {
		return fmt.Errorf(resRepoUpdate.Stderr)
	}

	taskUpgrade := execute.Task{
		Command:     helmPath,
		Args:        []string{
			"upgrade", "-i",
			"--kubeconfig", os.Getenv("KUBECONFIG"),
			"-f", "https://raw.githubusercontent.com/vertigobr/safira/master/k3d.yaml",
			"vtg-ipaas", "vtg-ipaas/vtg-ipaas",
		},
		StreamStdio: verboseFlag,
	}

	resUpgrade, err := taskUpgrade.Execute()
	if err != nil {
		return err
	}

	if resUpgrade.ExitCode != 0 {
		return fmt.Errorf(resUpgrade.Stderr)
	}

	return nil
}
