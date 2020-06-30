// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package config

import (
	"testing"
)

func TestCreateInBinDirDir(t *testing.T) {
	userDir, err := CreateInBinDir()

	if err != nil {
		t.Fatal("Não foi possível obter a pasta do usuário.")
	}

	t.Log(userDir)
}
