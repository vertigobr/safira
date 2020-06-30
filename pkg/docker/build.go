package docker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/vertigobr/safira/pkg/execute"
	"github.com/vertigobr/safira/pkg/utils"
)

type dockerBuild struct {
	path      string
	image     string
	noCache   bool
	buildArgs map[string]string
}

func Build(image, functionName, handler, language string, noCache bool, args map[string]string) error {
	if isValidTemplate(language) {
		buildPath, err := createBuildFolder(functionName, handler, language)
		if err != nil {
			return err
		}

		db := dockerBuild{
			path:      buildPath,
			image:     image,
			noCache:   noCache,
			buildArgs: args,
		}

		taskBuildArgs := getTaskBuildArgs(db)

		taskBuild := execute.Task{
			Command:      "docker",
			Args:         taskBuildArgs,
			StreamStdio:  true,
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
		return fmt.Errorf("atributo lang da função %s é inválido", functionName)
	}
}

func isValidTemplate(lang string) bool {
	if _, err := os.Stat("./template/" + lang); err == nil {
		return true
	}

	return false
}

func createBuildFolder(functionName, handler, language string) (string, error) {
	buildPath := fmt.Sprintf("./build/%s/", functionName)

	err := os.RemoveAll(buildPath)
	if err != nil {
		return "", fmt.Errorf("error ao limpar a pasta de build: %s\n", buildPath)
	}

	functionPath := path.Join(buildPath, "function")

	err = utils.Copy(path.Join("./template/", language), buildPath, false, true)
	if err != nil {
		return "", fmt.Errorf("error ao copiar arquivos para build: %s - %s.\n", buildPath, err.Error())
	}

	err = utils.Copy(filepath.Clean(handler), filepath.Clean(functionPath), true, true)
	if err != nil {
		return "", fmt.Errorf("error ao copiar arquivos para build: %s - %s.\n", buildPath, err.Error())
	}

	return buildPath, nil
}

func getTaskBuildArgs(db dockerBuild) (args []string) {
	args = append(args, "build")
	args = append(args, "--tag", db.image)

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
