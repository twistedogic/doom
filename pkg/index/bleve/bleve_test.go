package bleve

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/twistedogic/doom/testutil"
)

func Test_Index(t *testing.T) {
	dir, err := ioutil.TempDir("", "bleve_index")
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() { os.RemoveAll(dir) }
	file := filepath.Join(dir, "index")
	s := testutil.NewMockStore(t, make(map[string][]byte), false)
	index, err := New(file, s)
	if err != nil {
		t.Fatal(err)
	}
	testutil.IndexTest(t, index, cleanup)
}
