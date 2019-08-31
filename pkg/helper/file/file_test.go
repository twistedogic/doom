package file

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func SetupDir(t *testing.T) string {
	t.Helper()
	name := strconv.Itoa(rand.Int())
	dir, err := ioutil.TempDir("", name)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestCreateIfNotExist(t *testing.T) {

	cases := map[string]struct {
		filename string
		exist    bool
	}{
		"base":  {"test", false},
		"exist": {"test", true},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			dir := SetupDir(t)
			defer os.RemoveAll(dir)
			filename := filepath.Join(dir, tc.filename)
			if tc.exist {
				if _, err := os.Create(filename); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := CreateIfNotExist(filename); err != nil {
				t.Fatal(err)
			}
		})
	}
}
