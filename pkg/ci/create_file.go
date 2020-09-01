package ci

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/utils"
	"gopkg.in/gookit/color.v1"
	y "gopkg.in/yaml.v2"
)

func CreateFile() error {
	gitlabCi := GitlabCi{
		Image: "vertigo/safira:latest",
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
		return fmt.Errorf("%s Reading error to %s, check file", color.Red.Text("[!]"), GitlabCiFileName)
	}

	if err := utils.CreateYamlFile(GitlabCiFileName, yamlBytes, true); err != nil {
		return err
	}

	//if err := AppendFunction(functionName); err != nil {
	//	return err
	//}

	return nil
}

//func AppendFunction(functionName string) error {
//	deployName := fmt.Sprintf("%s:deploy", functionName)
//	undeployName := fmt.Sprintf("%s:undeploy", functionName)
//
//	jobs := FunctionsJobs{
//		Jobs: map[string]Job{
//			deployName: {
//				Name:  deployName,
//				Stage: "deploy",
//				Script: []string{
//					fmt.Sprintf("safira function deploy %s --kubeconfig=${KUBECONFIG}", functionName),
//				},
//			},
//			undeployName: {
//				Name:  undeployName,
//				Stage: "undeploy",
//				Script: []string{
//					fmt.Sprintf("safira function undeploy %s --kubeconfig=${KUBECONFIG}", functionName),
//				},
//			},
//		},
//	}
//
//	yamlBytes, err := y.Marshal(&jobs)
//	if err != nil {
//		return fmt.Errorf("error ao executar o marshal para o arquivo %s: %s", GitlabCiFileName, err.Error())
//	}
//
//	if err := utils.AppendYamlFile(GitlabCiFileName, yamlBytes); err != nil {
//		return err
//	}
//
//	return nil
//}
