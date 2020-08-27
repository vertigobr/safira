package stack

import (
	"fmt"
	"os"

	"github.com/vertigobr/safira/pkg/utils"
	"gopkg.in/gookit/color.v1"
	y "gopkg.in/yaml.v2"
)

func RemoveFunction(functionName string, removeFolder, verboseFlag bool) error {
	stack, err := LoadStackFile("")
	if err != nil {
		return err
	}

	if verboseFlag {
		fmt.Printf("%s Checking the existence of the function\n", color.Blue.Text("[v]"))
	}

	if _, exists := stack.Functions[functionName]; !exists {
		return fmt.Errorf("%s Function not found in the stack file", color.Red.Text("[!]"))
	}

	if verboseFlag {
		fmt.Printf("%s Removing function %s from project\n", color.Blue.Text("[v]"), functionName)
	}

	if removeFolder {
		functionPath := stack.Functions[functionName].Handler
		if err := os.RemoveAll(functionPath); err != nil {
			return fmt.Errorf("%s Error removing folder from function %s", color.Red.Text("[!]"), functionName)
		}
	}

	delete(stack.Functions, functionName)

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return fmt.Errorf("%s Error processing %s", color.Red.Text("[!]"), GetStackFileName())
	}

	if err := utils.CreateYamlFile(GetStackFileName(), yamlBytes, true); err != nil {
		return err
	}

	return nil
}
