package deploy

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"strings"
)

func createYamlFile(fileName string, bytes []byte) error {
	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		fileName = fileName + ".yml"
	}

	f, err := os.OpenFile(fileName, os.O_SYNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Truncate(0)

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func readFileEnv() error {
	if err := gotenv.Load(); err != nil {
		return err
	}

	return nil
}

func getProjectName() (string, error) {
	s := os.Getenv("PROJECT_NAME")
	if len(s) == 0 {
		return "", fmt.Errorf("variável PROJECT_NAME não encontrada no arquivo .env")
	}

	return s, nil
}

func getImageName() (string, error) {
	s := os.Getenv("IMAGE")
	if len(s) == 0 {
		return "", fmt.Errorf("variável IMAGE não encontrada no arquivo .env")
	}

	return s, nil
}

func getDomain() (string, error) {
	s := os.Getenv("DOMAIN")
	if len(s) == 0 {
		return "", fmt.Errorf("variável DOMAIN não encontrada no arquivo .env")
	}

	return s, nil
}

func getPort() (int, error) {
	s := os.Getenv("PORT")
	if len(s) == 0 {
		return 0, fmt.Errorf("variável PORT não encontrada no arquivo .env")
	}

	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return port, nil
}
