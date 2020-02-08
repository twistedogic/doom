package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/testutil"
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

type TestWrapper struct {
	t *testing.T
}

func (t TestWrapper) Setup() store.Store {
	t.t.Helper()
	fs := &mockFs{t.t, make(map[string][]byte)}
	return New(fs)
}

func (t TestWrapper) Cleanup() {}

func TestFileStore(t *testing.T) {
	testutil.StoreTest(t, TestWrapper{t})
}
