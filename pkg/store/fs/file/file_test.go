package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func setupTempFolder(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "testfs")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestWriteFile(t *testing.T) {
	content := "test"
	name := "newdir/doc"
	dir := setupTempFolder(t)
	defer os.RemoveAll(dir)
	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)
	writer := &Target{fs}
	buf := bytes.NewBufferString(content)
	if _, err := fs.Open(name); os.IsExist(err) {
		t.Fatal(err)
	}
	if err := writer.WriteFile(name, buf); err != nil {
		t.Fatal(err)
	}
	f, err := fs.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != content {
		t.Fatal(string(b))
	}
}
