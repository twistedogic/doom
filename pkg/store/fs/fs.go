package fs

import (
	"bytes"
	"io"
	"strings"
)

type FileSystem interface {
	WriteFile(string, io.ReadWriter) error
	ReadFile(string, io.Writer) error
	List() ([]string, error)
}

type FileStore struct {
	FileSystem
}

func New(fs FileSystem) FileStore {
	return FileStore{fs}
}

func (f FileStore) Get(key string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := f.ReadFile(key, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f FileStore) Set(key string, b []byte) error {
	buf := bytes.NewBuffer(b)
	return f.WriteFile(key, buf)
}

func (f FileStore) Scan(pattern ...string) ([]string, error) {
	list, err := f.List()
	if len(pattern) == 0 {
		return list, err
	}
	matches := make([]string, 0, len(list))
	p := pattern[0]
	for _, key := range list {
		if strings.Contains(key, p) {
			matches = append(matches, key)
		}
	}
	return matches, nil
}
