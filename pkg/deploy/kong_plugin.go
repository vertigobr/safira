package deploy

import y "gopkg.in/yaml.v2"

type kongPlugin struct {
	ApiVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   kongPluginMetadata `yaml:"metadata"`
	Config     kongPluginConfig   `yaml:"spec"`
	Plugin     string             `yaml:"plugin"`
}

type kongPluginMetadata struct {
	Name string `yaml:"name"`
}

type kongPluginConfig struct {
	Minute  int64  `yaml:"minute"`
	LimitBy string `yaml:"limit_by"`
	Policy  string `yaml:"policy"`
}

func CreateYamlKongPlugin(fileName string) error {
	kongPlugin := kongPlugin{
		ApiVersion: "configuration.konghq.com/v1",
		Kind:       "KongPlugin",
		Metadata: kongPluginMetadata{
			Name: "rate-limit-by-ip-5-min",
		},
		Config: kongPluginConfig{
			Minute: 5,
			LimitBy: "ip",
			Policy: "local",
		},
		Plugin: "rate-limiting",
	}

	yamlBytes, err := y.Marshal(&kongPlugin)
	if err != nil {
		return err
	}

	if err := createYamlFile(fileName, yamlBytes); err != nil {
		return err
	}

	return nil
}

