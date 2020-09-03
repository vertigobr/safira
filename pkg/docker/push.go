package docker

import (
	"fmt"

	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/git"
)

func Push(image string, useSha bool) error {
	taskImagePush := execute.Task{
		Command: "docker",
		Args: []string{
			"push", image,
		},
		StreamStdio: true,
	}

	res, err := taskImagePush.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}

	if useSha {
		imageWithCommitSha, _ := git.GetImageWithCommitSha(image)
		if len(imageWithCommitSha) > 0 {
			taskImageCommitShaPush := execute.Task{
				Command: "docker",
				Args: []string{
					"push", imageWithCommitSha,
				},
				StreamStdio: true,
			}

			res, err := taskImageCommitShaPush.Execute()
			if err != nil {
				return err
			}

			if res.ExitCode != 0 {
				return fmt.Errorf(res.Stderr)
			}
		}
	}

	return nil
}
