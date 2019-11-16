package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type Target struct {
	fs afero.Fs
}

func New() *Target {
	fs := afero.NewOsFs()
	return &Target{fs}
}

func (t *Target) CreateIfNotExist(name string) (afero.File, error) {
	if dir := filepath.Dir(name); dir != "" {
		if err := t.fs.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	f, err := t.fs.Create(name)
	if os.IsExist(err) {
		return t.fs.Open(name)
	}
	return f, err
}

func (t *Target) WriteFile(name string, r io.ReadWriter) error {
	f, err := t.CreateIfNotExist(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return err
	}
	return nil
}
