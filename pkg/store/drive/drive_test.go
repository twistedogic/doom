package drive

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testCred  = "../../../../testdata/cred.json"
	testCache = "../../../../testdata/cache.json"
)

func setup(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(testCred); err != nil {
		t.Skip()
		return
	}
	if _, err := os.Stat(testCache); err != nil {
		t.Skip()
		return
	}
}

func TestDrive(t *testing.T) {
	setup(t)
	cases := map[string]struct {
		filename string
		path     string
		content  []byte
	}{
		"base":              {"base", "base", []byte("base")},
		"with dir":          {"dir_base", "dir1/dir_base", []byte("dir_base")},
		"with multiple dir": {"basename", "dir1/dir2/basename", []byte("basename")},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			r := bytes.NewBuffer(tc.content)
			d, err := New(testCred, testCache)
			if err != nil {
				t.Fatal(err)
			}
			id, err := d.WriteFile(tc.path, r)
			if err != nil {
				t.Fatal(err)
			}
			buf := &bytes.Buffer{}
			if err := d.Download(id, buf); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(buf.String(), string(tc.content)); diff != "" {
				t.Fatal(diff)
			}
			if err := d.Delete(id); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMkdirAll(t *testing.T) {
	setup(t)
	folderName := fmt.Sprintf("testfolder-%d", time.Now().Unix())
	d, err := New(testCred, testCache)
	if err != nil {
		t.Fatal(err)
	}
	id, err := d.MkdirAll(folderName)
	if err != nil {
		t.Fatal(err)
	}
	oid, err := d.MkdirAll(folderName)
	if err != nil {
		t.Fatal(err)
	}
	if id != oid {
		t.Fatalf("want %s, got %s", id, oid)
	}
	if err := d.Delete(id); err != nil {
		t.Fatal(err)
	}
	if _, err := d.IsExist(folderName); err == nil {
		t.Fatal(err)
	}
}
