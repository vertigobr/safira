package stack

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

func AppendFunction(function Function) error {
	stack, err := LoadStackFile(stackFileName)
	if err != nil {
		return err
	}

	if stack.Functions[function.Name].Handler != "" {
		return fmt.Errorf(fmt.Sprintf("\nNome da função em uso, tente usar outro!"))
	}

	stack.Functions[function.Name] = function

	yamlBytes, err := y.Marshal(&stack)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", stackFileName, err.Error())
	}

	return utils.CreateYamlFile(stackFileName, yamlBytes, true)
}
