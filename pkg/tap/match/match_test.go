package match

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/twistedogic/doom/pkg/model"
)

const testdataPath = "../../../testdata"

func SetFile(t *testing.T) (string, func()) {
	t.Helper()
	f, err := ioutil.TempFile("", "store")
	if err != nil {
		t.Fatal(err)
	}
	return f.Name(), func() {
		f.Close()
		os.Remove(f.Name())
	}
}

func SetupServer(t *testing.T, dir string) *httptest.Server {
	t.Helper()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	cache := make(map[string][]byte)
	for _, info := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, info.Name()))
		if err != nil {
			t.Fatal(err)
		}
		cache[info.Name()] = b
	}
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		for name, b := range cache {
			if strings.Contains(req.URL.Path, name) {
				res.Write(b)
				return
			}
		}
	}))
}

func TestUpdate(t *testing.T) {
	ts := SetupServer(t, testdataPath)
	defer ts.Close()
	file, cleanup := SetFile(t)
	defer cleanup()
	app, err := New(file, ts.URL, -10, -1)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Backfill(); err != nil {
		t.Fatal(err)
	}
	var d model.Detail
	if err := app.GetLastest(&d); err != nil {
		t.Fatal(err)
	}
	var m model.Match
	if err := app.GetLastest(&m); err != nil {
		t.Fatal(err)
	}
	t.Log(m)
}
