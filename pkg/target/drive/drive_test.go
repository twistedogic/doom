package drive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testCred  = "../../../testdata/cred.json"
	testCache = "../../../testdata/cache.json"
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

type testItem struct {
	Name string
}

func TestDrive(t *testing.T) {
	setup(t)
	cases := map[string]struct {
		filename string
		path     string
		content  testItem
	}{
		"base":              {"base", "base", testItem{"base"}},
		"with dir":          {"dir_base", "dir1/dir_base", testItem{"dir_base"}},
		"with multiple dir": {"basename", "dir1/dir2/basename", testItem{"basename"}},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			d, err := New(tc.path, testCred, testCache)
			if err != nil {
				t.Fatal(err)
			}
			if err := d.UpsertItem(tc.content); err != nil {
				t.Fatal(err)
			}
			if err := d.Close(); err != nil {
				t.Fatal(err)
			}
			id, exist, err := d.IsExist(tc.filename)
			if err != nil || !exist {
				t.Fatal(err)
			}
			buf := &bytes.Buffer{}
			if err := d.Download(id, buf); err != nil {
				t.Fatal(err)
			}
			var item testItem
			if err := json.NewDecoder(buf).Decode(&item); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(testItem{tc.filename}, item); diff != "" {
				t.Fatal(diff)
			}
			if err := d.Delete(id); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateFolders(t *testing.T) {
	setup(t)
	folderName := fmt.Sprintf("testfolder-%d", time.Now().Unix())
	d, err := New("", testCred, testCache)
	if err != nil {
		t.Fatal(err)
	}
	id, err := d.createFolders(folderName)
	if err != nil {
		t.Fatal(err)
	}
	oid, err := d.createFolders(folderName)
	if err != nil {
		t.Fatal(err)
	}
	if id != oid {
		t.Fatalf("want %s, got %s", id, oid)
	}
	if err := d.Delete(id); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := d.IsExist(folderName); err != nil || ok {
		t.Fatal(err)
	}
}
