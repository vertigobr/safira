package deploy

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/utils"
	y "gopkg.in/yaml.v2"
)

// Yaml file structure
type K8sYaml struct {
	ApiVersion string      `yaml:"apiVersion,omitempty"`
	Kind       string      `yaml:"kind,omitempty"`
	Metadata   metadata    `yaml:"metadata,omitempty"`
	Spec       interface{} `yaml:"spec,omitempty"`
}

type metadata struct {
	Name   string         `yaml:"name,omitempty"`
	Labels metadataLabels `yaml:"labels,omitempty"`
}

type metadataLabels struct {
	App string `yaml:"app,omitempty"`
}

func (k *K8sYaml) CreateYamlFile(fileName string) error {
	yamlBytes, err := y.Marshal(&k)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", fileName, err.Error())
	}

	if err := utils.CreateYamlFile(fileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}
