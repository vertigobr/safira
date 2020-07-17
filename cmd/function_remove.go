// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/stack"

	"github.com/spf13/cobra"
)

var functionRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Removes the function from the project",
	Long:    "Removes the function from the project",
	Example: `To remove the function from the project, run:

    $ safira function undeploy function-name`,
	RunE:    runFunctionRemove,
}

func init() {
	functionCmd.AddCommand(functionRemoveCmd)
	functionRemoveCmd.Flags().BoolP("remove-folder", "R", false, "remove folder from function")
}

func runFunctionRemove(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	removeFolderFlag, _ := cmd.Flags().GetBool("remove-folder")

	if err := stack.RemoveFunction(args[0], removeFolderFlag, verboseFlag); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Function %s successfully removed!", args[0]))
	return nil
}
