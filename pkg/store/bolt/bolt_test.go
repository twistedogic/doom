package bolt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/testutil"
)

type TestWrapper struct {
	*testing.T
	dir string
}

func Setup(t *testing.T) *TestWrapper {
	return &TestWrapper{t, ""}
}

func (t *TestWrapper) Setup() store.Store {
	t.Helper()
	dir, err := ioutil.TempDir("", "tmpdb")
	if err != nil {
		t.Fatal(err)
	}
	t.dir = dir
	file := filepath.Join(dir, "tempdb")
	s, err := New(file)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func (t *TestWrapper) Cleanup() {
	os.RemoveAll(t.dir)
}

func Test_Store(t *testing.T) {
	testutil.StoreTest(t, Setup(t))
}
