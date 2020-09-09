// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	s "github.com/vertigobr/safira/pkg/stack"
	"io/ioutil"
)

func (k *K8sYaml) MountConfigMap(configMapName, swaggerFile, repoName, env string) error {
	b, _ := ioutil.ReadFile(swaggerFile)

	stack, err := s.LoadStackFile(env)
	if err != nil {
		return err
	}

	configMapName = GetDeployName(stack, configMapName)

	*k = K8sYaml{
		ApiVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: metadata{
			Name: configMapName,
			Labels: map[string]string{
				"name": configMapName,
			},
			Annotations: map[string]string{
				"safira.io/repository": repoName,
			},
		},
		Data: map[string]string{
			"swagger": string(b),
		},
	}

	return nil
}
