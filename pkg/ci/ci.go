// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package ci

const gitlabCiFileName = ".gitlab-ci.yml"

type Job struct {
	Name   string   `yaml:"-"`
	Stage  string   `yaml:"stage"`
	Script []string `yaml:"script"`
}

type GitlabCi struct {
	Image        string   `yaml:"image,omitempty"`
	Services     []string `yaml:"services,omitempty"`
	Stages       []string `yaml:"stages,omitempty"`
	BeforeScript []string `yaml:"before_script,omitempty"`
	Publish      Job      `yaml:"publish,omitempty"`
	Deploy       Job      `yaml:"deploy,omitempty"`
}


