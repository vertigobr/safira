// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/gookit/color.v1"
)

// Used as a reference: https://github.com/alexellis/arkade/blob/master/pkg/helm/untar.go
func Untar(r io.Reader, dir string) (err error) {
	t0 := time.Now()
	nFiles := 0

	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("%s Require gzip-compressed", color.Red.Text("[!]"))
	}

	tr := tar.NewReader(zr)
	loggedChtimesError := false

	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}

		if strings.Contains(f.Name, "LICENSE") {
			break
		}

		if strings.Contains(f.Name, "README") {
			break
		}

		if err != nil {
			return fmt.Errorf("%s Error reading %s", color.Red.Text("[!]"), dir)
		}

		if !validRelPath(f.Name) {
			return fmt.Errorf("%s Invalid name: %q", color.Red.Text("[!]"), f.Name)
		}

		baseFile := filepath.Base(f.Name)
		abs := path.Join(dir, baseFile)

		fi := f.FileInfo()
		mode := fi.Mode()

		switch {
		case mode.IsDir():
			break

		case mode.IsRegular():
			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}

			n, err := io.Copy(wf, tr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}

			if err != nil {
				return fmt.Errorf("%s Error writing in %s", color.Red.Text("[!]"), abs)
			}

			if n != f.Size {
				return fmt.Errorf("%s %d bytes %s; excede %d", color.Red.Text("[!]"), n, abs, f.Size)
			}

			modTime := f.ModTime
			if modTime.After(t0) {
				modTime = t0
			}

			if !modTime.IsZero() {
				if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
					log.Printf("%s Error when changing modtime: %v", color.Red.Text("[!]"), err)
					loggedChtimesError = true
				}
			}

			nFiles++
		default:
		}
	}

	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}

	return true
}
