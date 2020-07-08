// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	k "github.com/vertigobr/safira/pkg/kong"
)

var kongNewCmd = &cobra.Command{
	Use:     "new [TYPE] [NAME]",
	Short:   "Create a yaml structure for the Kong CRD",
	Long:    "Create a yaml structure for the Kong CRD",
	Example: `To create a yaml file for kong, run:

    $ safira kong new plugin plugin-name`,
	PreRunE: preRunKongpluginNew,
	RunE:    runKongpluginNew,
}

func init() {
	kongCmd.AddCommand(kongNewCmd)
}

func preRunKongpluginNew(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		_ = cmd.Help()
		os.Exit(0)
	} else if len(args) == 1 {
		return fmt.Errorf("nome/Tipo do recurso não inserido")
	}

	return nil
}

func runKongpluginNew(cmd *cobra.Command, args []string) error {
	assetType := args[0]
	assetName := args[1]

	fmt.Println(fmt.Sprintf("[+] Criando o plugin %s do tipo %s", assetName, assetType))
	if err := k.Create(assetName, assetType); err != nil {
		return err
	}

	fmt.Println("\nPlugin criado com sucesso!")

	return nil
}
