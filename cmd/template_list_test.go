// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"testing"

	"github.com/vertigobr/safira/pkg/execute"
)

func TestTemplateList(t *testing.T) {
	taskTest := execute.Task{
		Command:     "safira",
		Args:        []string{
			"template", "list",
		},
		StreamStdio:  true,
		PrintCommand: true,
	}

	t.Log("Iniciando execução: safira template list")
	res, err := taskTest.Execute()
	if err != nil {
		t.Errorf(err.Error())
	}

	if res.ExitCode != 0 {
		t.Errorf(res.Stderr)
	}
}