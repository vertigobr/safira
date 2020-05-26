package deploy

import y "gopkg.in/yaml.v2"

type ingress struct {
	ApiVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   ingressMetadata `yaml:"metadata"`
	Spec       ingressSpec     `yaml:"spec"`
}

type ingressMetadata struct {
	Name string `yaml:"name"`
}

type ingressSpec struct {
	Rules []rule `yaml:"rules"`
}

type rule struct {
	Host string `yaml:"host"`
	Http http   `yaml:"http"`
}

type http struct {
	Paths []path `yaml:"paths"`
}

type path struct {
	Path    string  `yaml:"path"`
	Backend backend `yaml:"backend"`
}

type backend struct {
	ServiceName string `yaml:"serviceName"`
	ServicePort int    `yaml:"servicePort"`
}

func CreateYamlIngress(fileName string) error {
	if err := readFileEnv(); err != nil {
		return err
	}

	projectName, domain, port, err := getIngressEnvs()
	if err != nil {
		return nil
	}

	ingress := ingress{
		ApiVersion: "extensions/v1beta1",
		Kind:       "Ingress",
		Metadata: ingressMetadata{
			Name: projectName,
		},
		Spec: ingressSpec{
			Rules: []rule{
				{
					Host: domain,
					Http: http{
						Paths: []path{
							{
								Path: "/function/" + projectName,
								Backend: backend{
									ServiceName: projectName,
									ServicePort: port,
								},
							},
						},
					},
				},
			},
		},
	}

	yamlBytes, err := y.Marshal(&ingress)
	if err != nil {
		return err
	}

	if err := createYamlFile(fileName, yamlBytes); err != nil {
		return err
	}

	return nil
}

func getIngressEnvs() (string, string, int, error) {
	projectName, err := getProjectName()
	if err != nil {
		return "", "", 0, err
	}

	domain, err := getDomain()
	if err != nil {
		return "", "", 0, err
	}

	port, err := getPort()
	if err != nil {
		return "", "", 0, err
	}

	return projectName, domain, port, nil
}
