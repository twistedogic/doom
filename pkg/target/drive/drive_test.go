package drive

import (
	"fmt"
	"os"
	"testing"
	"time"
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
	}{
		"base":              {"test"},
		"with dir":          {"dir1/test"},
		"with multiple dir": {"dir1/dir2/test"},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			d, err := New(tc.filename, testCred, testCache)
			if err != nil {
				t.Fatal(err)
			}
			if err := d.UpsertItem(testItem{tc.filename}); err != nil {
				t.Fatal(err)
			}
			if err := d.Close(); err != nil {
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
