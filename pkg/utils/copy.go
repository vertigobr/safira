// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Copy(source, dest string, folders, files bool) error {
	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("error ao procurar a pasta template no repositório")
	}

	if !sourceInfo.IsDir() {
		return fmt.Errorf("error template não é um diretório")
	}

	_, err = os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error ao acessar pasta template de destino")
	}

	if err == nil {
		return fmt.Errorf("pasta template de destino não existe")
	}

	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("error ao criar pasta template de destino")
	}

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return fmt.Errorf("error ao ler pasta template do repositório")
	}

	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		dstPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if folders {
				err = Copy(srcPath, dstPath, folders, files)
				if err != nil {
					return err
				}
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(source string, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("error ao abrir o arquivo %s", source)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error ao criar o arquivo %s", dest)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("error ao copiar conteúdo do arquivo %s para %s", source, dest)
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("arquivo não existe: %s", source)
	}
	err = os.Chmod(dest, si.Mode())
	if err != nil {
		return fmt.Errorf("error ao aplicar chmod %s no arquivo %v", dest, si.Mode())
	}

	return nil
}

