package store

import (
  "bytes"
  "strings"
)

type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Scan(...string) ([]string, error)
}

func Transform(src, dst Store) error {
	keys, err := src.Scan()
	if err != nil {
		return err
	}
	for _, key := range keys {
		b, err := src.Get(key)
		if err != nil {
			return err
		}
		if err := dst.Set(key, b); err != nil {
			return err
		}
	}
	return nil
}

type FileSystem interface {
  WriteFile(string, io.ReadWriter) error
  ReadFile(string, io.Writer) error
  List() ([]string, error)
}

type FileStore struct {
  FileSystem
}

func NewFileStore(fs FileSystem) (*FileStore,error) {
  keys, err := fs.List()
  return &FileStore{keys, fs}, err
}

func (f *FileStore) Get(key string) ([]byte, error) {
  buf := new(bytes.Buffer)
  if err := f.ReadFile(key, buf); err != nil {
    return nil, err
  }
  return buf.Byte(), nil
}

func (f *FileStore) Set(key string, b []byte) error {
  buf := bytes.NewBuffer(b)
  return f.WriteFile(key, buf)
}

func (f *FileStore) Scan(pattern ...string) ([]string, error){
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

