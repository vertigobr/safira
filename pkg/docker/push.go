package docker

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/execute"
)

func Push(image string) error {
	taskPush := execute.Task{
		Command:      "docker",
		Args:         []string{
			"push", image,
		},
		StreamStdio:  true,
	}

	res, err := taskPush.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	return nil
}
