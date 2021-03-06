// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/stack"
	"gopkg.in/gookit/color.v1"
)

var functionRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the function from the project",
	Long:  "Removes the function from the project",
	Example: `To remove the function from the project, run:

    $ safira function remove function-name`,
	RunE:                       runFunctionRemove,
	SuggestionsMinimumDistance: 1,
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

	fmt.Printf("%s Function %s successfully removed\n", color.Red.Text("[!]"), args[0])
	return nil
}
