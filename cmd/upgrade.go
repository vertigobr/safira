// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
	"gopkg.in/gookit/color.v1"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "CLI updates itself",
	Long:  "CLI updates itself",
	Example: `To upgrade CLI, run:

    $ safira upgrade`,
	RunE:                       runUpgrade,
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringP("version", "v", "latest", "CLI download version")
}

func runUpgrade(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	if err := upgrade(verboseFlag); err != nil {
		return err
	}

	return nil
}

func upgrade(verboseFlag bool) error {
	fmt.Printf("%s Upgrading Safira\n", color.Green.Text("[+]"))

	tag, err := get.DownloadSafira(rootCmd.Version, verboseFlag)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s Upgrading version %s finished successfully\n", color.Cyan.Text("[✓]"), tag)

	return nil
}
