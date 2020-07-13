// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import (
	s "github.com/vertigobr/safira/pkg/stack"
)

func (k *K8sYaml) MountKongPlugin(functionName, pluginName string) error {
	stack, err := s.LoadStackFile()
	if err != nil {
		return err
	}

	kind := getKingKongPlugin(stack.Functions[functionName].Plugins[pluginName].Type)
	metadataName := functionName + "-" + pluginName

	*k = K8sYaml{
		ApiVersion: "configuration.konghq.com/v1",
		Kind:       kind,
		Metadata: metadata{
			Name: metadataName,
			Labels: map[string]string{
				"global": stack.Functions[functionName].Plugins[pluginName].Global,
			},
		},
		Config: stack.Functions[functionName].Plugins[pluginName].Config,
		ConfigFrom: stack.Functions[functionName].Plugins[pluginName].ConfigFrom,
		Plugin: pluginName,
	}

	return nil
}

func getKingKongPlugin(pluginType string) (kind string) {
	if pluginType == "cluster" || pluginType == "cluster-plugin" {
		kind = "KongClusterPlugin"
	} else {
		kind = "KongPlugin"
	}

	return
}
