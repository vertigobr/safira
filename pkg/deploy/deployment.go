// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

type deploymentSpec struct {
	Replicas int                    `yaml:"replicas,omitempty"`
	Selector deploymentSpecSelector `yaml:"selector,omitempty"`
	Template deploymentSpecTemplate `yaml:"template,omitempty"`
}

type deploymentSpecSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels,omitempty"`
}

type deploymentSpecTemplate struct {
	Metadata     metadata     `yaml:"metadata,omitempty"`
	TemplateSpec templateSpec `yaml:"spec,omitempty"`
}

type templateSpec struct {
	Containers []containers `yaml:"containers,omitempty"`
	Volumes    []volumes    `yaml:"volumes,omitempty"`
}

type containers struct {
	Name         string           `yaml:"name,omitempty"`
	Image        string           `yaml:"image,omitempty"`
	Ports        []containerPorts `yaml:"ports,omitempty"`
	Env          []containerEnv   `yaml:"env,omitempty"`
	VolumeMounts []volumeMounts   `yaml:"volumeMounts,omitempty"`
}

type containerPorts struct {
	ContainerPort int `yaml:"containerPort,omitempty"`
}

type containerEnv struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type volumeMounts struct {
	Name      string `yaml:"name,omitempty"`
	MountPath string `yaml:"mountPath,omitempty"`
}

type configMapItem struct {
	Key  string `yaml:"key,omitempty"`
	Path string `yaml:"path,omitempty"`
}
type volumes struct {
	Name      string          `yaml:"name,omitempty"`
	ConfigMap volumeConfigMap `yaml:"configMap,omitempty"`
}

type volumeConfigMap struct {
	Name  string          `yaml:"name,omitempty"`
	Items []configMapItem `yaml:"items,omitempty"`
}

func (k *K8sYaml) MountDeployment(deploymentName, imageName, path, repoName string) error {
	*k = K8sYaml{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: metadata{
			Name: deploymentName,
			Labels: map[string]string{
				"app": deploymentName,
			},
			Annotations: map[string]string{
				"safira.io/repository": repoName,
			},
		},
		Spec: deploymentSpec{
			Replicas: 1,
			Selector: deploymentSpecSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: deploymentSpecTemplate{
				Metadata: metadata{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				TemplateSpec: templateSpec{
					Containers: []containers{
						{
							Name:  deploymentName,
							Image: imageName,
							Ports: []containerPorts{
								{
									ContainerPort: 8080,
								},
							},
							Env: []containerEnv{
								{
									Name:  "BASE_URL",
									Value: path,
								},
								{
									Name:  "SWAGGER_JSON",
									Value: "/swagger-ui/swagger.yml",
								},
							},
							VolumeMounts: []volumeMounts{
								{
									Name:      deploymentName,
									MountPath: "/swagger-ui",
								},
							},
						},
					},
					Volumes: []volumes{
						{
							Name: deploymentName,
							ConfigMap: volumeConfigMap{
								Name: deploymentName,
								Items: []configMapItem{
									{
										Key:  "swagger",
										Path: "swagger.yml",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return nil
}
