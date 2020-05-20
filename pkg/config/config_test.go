package config

import "testing"

func TestCreateInBinDirDir(t *testing.T) {
	userDir, err := CreateInBinDir()

	if err != nil {
		t.Fatal("Não foi possível obter a pasta do usuário.")
	}

	t.Log(userDir)
}
