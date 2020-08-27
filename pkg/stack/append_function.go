package stack

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/utils"
	"gopkg.in/gookit/color.v1"
	y "gopkg.in/yaml.v2"
)

func AppendFunction(function Function) error {
	stack, err := LoadStackFile("")
	if err != nil {
		return err
	}

	if stack.Functions[function.Name].Handler != "" {
		return fmt.Errorf("%s Function name in use, try using another", color.Red.Text("[!]"))
	}

	stack.Functions[function.Name] = function

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return fmt.Errorf("%s Error reading %s file", color.Red.Text("[!]"), GetStackFileName())
	}

	return utils.CreateYamlFile(GetStackFileName(), yamlBytes, true)
}
