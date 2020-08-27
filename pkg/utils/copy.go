// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/gookit/color.v1"
)

func Copy(source, dest string, folders, files bool) error {
	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	sourceInfo, err := os.Stat(source)
	if err != nil || !sourceInfo.IsDir() {
		return fmt.Errorf("%s Could not find the template folder in the repository", color.Red.Text("[!]"))
	}

	_, err = os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%s Error accessing template folder", color.Red.Text("[!]"))
	}

	if err == nil {
		return fmt.Errorf("%s Template folder could not be found", color.Red.Text("[!]"))
	}

	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("%s Error creating template folder", color.Red.Text("[!]"))
	}

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return fmt.Errorf("%s Error reading contents of the repository template folder", color.Red.Text("[!]"))
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
		return fmt.Errorf("%s Error opening the %s file", color.Red.Text("[!]"), source)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("%s Error reading the %s file", color.Red.Text("[!]"), dest)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%s Error copying contents of file %s to %s", color.Red.Text("[!]"), source, dest)
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("%s %s file does not exist", color.Red.Text("[!]"), source)
	}
	err = os.Chmod(dest, si.Mode())
	if err != nil {
		return fmt.Errorf("%s Error changing %s file permission", color.Red.Text("[!]"), dest)
	}

	return nil
}
