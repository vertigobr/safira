// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	"fmt"

	s "github.com/vertigobr/safira/pkg/stack"
	"github.com/vertigobr/safira/pkg/utils"
)

func (k *K8sYaml) MountKongPlugin(functionName, pluginName, namespace, env string) error {
	stack, err := s.LoadStackFile(env)
	if err != nil {
		return err
	}

	metadataName := fmt.Sprintf("%s-%s", functionName, pluginName)

	repoName, err := utils.GetCurrentFolder()
	if err != nil {
		return err
	}

	*k = K8sYaml{
		ApiVersion: "configuration.konghq.com/v1",
		Kind:       "KongPlugin",
		Metadata: metadata{
			Name:      metadataName,
			Namespace: namespace,
			Labels: map[string]string{
				"global": stack.Functions[functionName].Plugins[pluginName].Global,
			},
			Annotations: map[string]string{
				"safira.io/repository": repoName,
				"safira.io/function":   functionName,
			},
		},
		Config:     stack.Functions[functionName].Plugins[pluginName].Config,
		ConfigFrom: stack.Functions[functionName].Plugins[pluginName].ConfigFrom,
		Plugin:     pluginName,
	}

	return nil
}
