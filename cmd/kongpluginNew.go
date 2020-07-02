// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	kp "github.com/vertigobr/safira/pkg/kongplugin"
)

var kongpluginNewCmd = &cobra.Command{
	Use:     "new [PLUGIN_NAME] --type [PLUGIN_TYPE]",
	Short:   "Cria a estrutura de um plugin",
	Long:    "Cria a estrutura de um plugin",
	PreRunE: preRunKongpluginNew,
	RunE:    runKongpluginNew,
}

func init() {
	kongpluginCmd.AddCommand(kongpluginNewCmd)
	kongpluginNewCmd.Flags().StringP("type", "t", "plugin", "Tipo de plugin que será gerado")
}

func preRunKongpluginNew(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nome do plugin não inserido")
	}

	return nil
}

func runKongpluginNew(cmd *cobra.Command, args []string) error {
	typeFlag, _ := cmd.Flags().GetString("type")

	for _, pluginName := range args {
		if err := kp.CreatePlugin(pluginName, typeFlag); err != nil {
			return err
		}
	}

	return nil
}
