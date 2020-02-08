package file

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type Target struct {
	fs afero.Fs
}

func New(fs afero.Fs) Target {
	return Target{fs}
}

func (t Target) CreateIfNotExist(name string) (afero.File, error) {
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

func (t Target) WriteFile(name string, r io.ReadWriter) error {
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

func (t Target) ReadFile(name string, w io.Writer) error {
	b, err := afero.ReadFile(t.fs, name)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (t Target) List() ([]string, error) {
	infos, err := afero.ReadDir(t.fs, ".")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(infos))
	for _, info := range infos {
		if !info.IsDir() {
			names = append(names, info.Name())
		}
	}
	return names, nil
}
