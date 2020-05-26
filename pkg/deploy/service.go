package deploy

import (
	y "gopkg.in/yaml.v2"
)

type service struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   serviceMetadata `yaml:"metadata"`
	Spec       serviceSpec     `yaml:"spec"`
}

type serviceMetadata struct {
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}

type serviceSpec struct {
	Type         string            `yaml:"type"`
	ExternalName string            `yaml:"externalName"`
	Ports        []port            `yaml:"ports"`
}

type port struct {
	Port int `yaml:"port"`
}


func CreateYamlService(fileName string) error {
	if err := readFileEnv(); err != nil {
		return err
	}

	projectName, p, err := getServiceEnvs()
	if err != nil {
		return nil
	}

	service := service{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: serviceMetadata{
			Name:   projectName,
			Labels: map[string]string{
				"app": projectName,
			},
			Annotations: map[string]string{
				"konghq.com/plugins": "prometheus",
			},
		},
		Spec: serviceSpec{
			Type: "ExternalName",
			ExternalName: "gateway",
			Ports: []port{
				{
					Port: p,
				},
			},
		},
	}

	yamlBytes, err := y.Marshal(&service)
	if err != nil {
		return err
	}

	if err := createYamlFile(fileName, yamlBytes); err != nil {
		return err
	}

	return nil
}

func getServiceEnvs() (string, int, error) {
	projectName, err := getProjectName()
	if err != nil {
		return "", 0, err
	}

	port, err := getPort()
	if err != nil {
		return "", 0, err
	}

	return projectName, port, nil
}
