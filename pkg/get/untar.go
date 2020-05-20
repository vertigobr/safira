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
)

// Usado como base: https://github.com/alexellis/arkade/blob/master/pkg/helm/untar.go

func Untar(r io.Reader, dir string) (err error) {
	t0 := time.Now()
	nFiles := 0

	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("requer gzip-compressed: %v", err)
	}

	tr := tar.NewReader(zr)
	loggedChtimesError := false

	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("falha na leitura do tar: %v", err)
			return fmt.Errorf("error no tar: %v", err)
		}

		if !validRelPath(f.Name) {
			return fmt.Errorf("tar contêm nome inválido: %q", f.Name)
		}

		baseFile := filepath.Base(f.Name)
		abs := path.Join(dir, baseFile)
		//fmt.Println(abs, f.Name)

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
				return fmt.Errorf("error ao escrever %s: %v", abs, err)
			}

			if n != f.Size {
				return fmt.Errorf("%d bytes %s; excede %d", n, abs, f.Size)
			}

			modTime := f.ModTime
			if modTime.After(t0) {
				modTime = t0
			}

			if !modTime.IsZero() {
				if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
					log.Printf("error ao alterar o modtime: %v", err)
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