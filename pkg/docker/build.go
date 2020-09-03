package docker

import (
	"fmt"
	"github.com/vertigobr/safira/pkg/git"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/utils"
	"gopkg.in/gookit/color.v1"
)

type dockerBuild struct {
	path      string
	image     string
	noCache   bool
	buildArgs map[string]string
}

func Build(image, functionName, handler, language string, useSha, noCache bool, args map[string]string, verboseFlag bool) error {
	if isValidTemplate(language) {
		commitSha := ""
		if useSha {
			commitSha, _ = git.GetCommitSha()
		}

		buildPath, err := createBuildFolder(functionName, handler, language, verboseFlag)
		if err != nil {
			return err
		}

		db := dockerBuild{
			path:      buildPath,
			image:     image,
			noCache:   noCache,
			buildArgs: args,
		}

		taskBuildArgs := getTaskBuildArgs(db, commitSha)

		taskBuild := execute.Task{
			Command:     "docker",
			Args:        taskBuildArgs,
			StreamStdio: true,
		}

		if verboseFlag {
			fmt.Printf("%s Running build\n", color.Blue.Text("[v]"))
		}

		res, err := taskBuild.Execute()
		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		return nil
	} else {
		return fmt.Errorf(
			"%s The input from %s to lang is invalid, look at templates running: %s",
			color.Red.Text("[!]"), language, color.Bold.Text("safira template list"),
		)
	}
}

func isValidTemplate(lang string) bool {
	if _, err := os.Stat("./template/" + lang); err == nil {
		return true
	}

	return false
}

func createBuildFolder(functionName, handler, language string, verboseFlag bool) (string, error) {
	buildPath := fmt.Sprintf("./build/%s/", functionName)

	err := os.RemoveAll(buildPath)
	if err != nil {
		return "", fmt.Errorf("%s Error cleaning build folder: %s", color.Red.Text("[!]"), buildPath)
	}

	if verboseFlag {
		fmt.Printf("%s Build %s function folder removed\n", color.Blue.Text("[v]"), functionName)
	}

	functionPath := path.Join(buildPath, "function")

	if verboseFlag {
		fmt.Printf("%s Copying build artifacts\n", color.Blue.Text("[v]"))
	}

	err = utils.Copy(path.Join("./template/", language), buildPath, false, true)
	if err != nil {
		return "", fmt.Errorf("%s Error copying template %s files to the build folder: %s", color.Red.Text("[!]"), language, buildPath)
	}

	if verboseFlag {
		fmt.Printf("%s Copying files from function %s\n", color.Blue.Text("[v]"), functionName)
	}

	err = utils.Copy(filepath.Clean(handler), filepath.Clean(functionPath), true, true)
	if err != nil {
		return "", fmt.Errorf("%s Error copying function %s files to the build folder: %s", color.Red.Text("[!]"), functionName, buildPath)
	}

	return buildPath, nil
}

func getTaskBuildArgs(db dockerBuild, commitSha string) (args []string) {
	args = append(args, "build")
	args = append(args, "--tag", db.image)

	if len(commitSha) > 0 {
		untaggedImage := strings.Split(db.image, ":")

		args = append(args, "--tag", strings.Replace(db.image, untaggedImage[len(untaggedImage)-1], commitSha, 1))
	}

	if db.noCache {
		args = append(args, "--no-cache")
	}

	for index, arg := range db.buildArgs {
		keyVal := index + "=" + arg
		args = append(args, "--build-arg", keyVal)
	}

	args = append(args, db.path)
	return
}
