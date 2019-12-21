package store

import (
  "fmt"
  "io/ioutil"
  "testing"
)

type mockFs struct {
  t *testing.T
  content map[string][]byte
}

func (m *mockFs) WriteFile(key string, r io.ReadWriter) error {
  b, err := ioutil.ReadAll(r)
  if err != nil {
    return err
  }
  m.content[key] = b
  return nil
}

func (m *mockFs)  ReadFile(key string, w io.Writer) error {
  b, ok := m.content[key]
  if !ok {
    return fmt.Errorf("no key found")
  }
  _, err := w.Write(b)
  return err
}

func (m *mockFs)  List() ([]string, error) {
  keys := make([]string,0,len(m.content))
  for k := range m.content {
    keys = append(keys, k)
  }
  return keys, nil
}
