package stack

import (
	"fmt"
	"os"

	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func RemoveFunction(functionName string, verboseFlag bool) error {
	stack, err := LoadStackFile()
	if err != nil {
		return err
	}

	if verboseFlag {
		fmt.Println("[+] Checking existence of the function")
	}

	if _, exists := stack.Functions[functionName]; !exists {
		return fmt.Errorf("function not found in the stack")
	}

	if verboseFlag {
		fmt.Println("[+] Removing function from project")
	}

	functionPath := stack.Functions[functionName].Handler
	delete(stack.Functions, functionName)

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return fmt.Errorf("error processing %s: %s", GetYamlFileName(), err.Error())
	}

	if err := utils.CreateYamlFile(GetYamlFileName(), yamlBytes, true); err != nil {
		return err
	}

	if err := os.RemoveAll(functionPath); err != nil {
		return fmt.Errorf("error removing folder from function %s", functionName)
	}

	return nil
}
