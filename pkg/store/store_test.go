package store

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockFs struct {
	t       *testing.T
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

func (m *mockFs) ReadFile(key string, w io.Writer) error {
	b, ok := m.content[key]
	if !ok {
		return fmt.Errorf("no key found")
	}
	_, err := w.Write(b)
	return err
}

func (m *mockFs) List() ([]string, error) {
	keys := make([]string, 0, len(m.content))
	for k := range m.content {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestFileStore(t *testing.T) {
	cases := map[string]struct {
		input map[string][]byte
	}{
		"base": {
			input: map[string][]byte{
				"id1": []byte(`ok`),
				"id2": []byte(`ok`),
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			fs := &mockFs{t, make(map[string][]byte)}
			store := NewFileStore(fs)
			for k, b := range tc.input {
				if err := store.Set(k, b); err != nil {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(tc.input, fs.content); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
