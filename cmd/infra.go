// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/get"
)

var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Manage the local infrastructure",
	Long:  "Manage the local infrastructure",
	SuggestionsMinimumDistance: 1,
}

func init() {
	rootCmd.AddCommand(infraCmd)
}

func checkInfra(verboseFlag bool) error {
	exist, err := get.CheckBinary(kubectlBinaryName, false, verboseFlag)
	if err != nil {
		return err
	} else if !exist {
		return fmt.Errorf(notExistBinary)
	}

	exist, err = get.CheckBinary(helmBinaryName, false, verboseFlag)
	if err != nil {
		return err
	} else if !exist {
		return fmt.Errorf(notExistBinary)
	}

	exist, err = get.CheckBinary(k3dBinaryName, false, verboseFlag)
	if err != nil {
		return err
	} else if !exist {
		return fmt.Errorf(notExistBinary)
	}

	return nil
}
