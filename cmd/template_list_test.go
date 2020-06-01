package cmd

import (
	"github.com/vertigobr/safira/pkg/execute"
	"testing"
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