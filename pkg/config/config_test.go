package config

import (
	"testing"
)

func TestCreateDir(t *testing.T) {
	userDir, err := CreateInBinDir("/.helm/")

	if err != nil {
		t.Fatal("Não foi possível obter a pasta do usuário.")
	}

	t.Log(userDir)
}
