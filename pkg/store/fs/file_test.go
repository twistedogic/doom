package fs

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestWriteFile(t *testing.T) {
	content := "test"
	name := "newdir/doc"
	fs := afero.NewMemMapFs()
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
