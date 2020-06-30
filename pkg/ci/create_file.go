package ci

import (
	"fmt"

	y "gopkg.in/yaml.v2"

	"github.com/vertigobr/safira/pkg/utils"
)

func CreateFile() error {
	gitlabCi := GitlabCi{
		Image:   "vertigobr/safira:latest",
		Services: []string{
			"docker:19.03.8-dind",
		},
		Stages: []string{
			"publish",
			"deploy",
		},
		BeforeScript: []string{
			"safira init",
		},
		Publish: Job{
			Stage: "publish",
			Script: []string{
				"safira template pull",
				"echo ${PASSWORD} | docker login -u ${USER} --password-stdin",
				"safira function build-push -A",
			},
		},
		Deploy: Job{
			Stage: "deploy",
			Script: []string{
				"safira function deploy -A --kubeconfig=${KUBECONFIG}",
			},
		},
	}

	yamlBytes, err := y.Marshal(&gitlabCi)
	if err != nil {
		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", gitlabCiFileName, err.Error())
	}

	if err := utils.CreateYamlFile(gitlabCiFileName, yamlBytes, true); err != nil {
		return err
	}

	return nil
}

