package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/twistedogic/doom/pkg/model"
)

const testdataPath = "../../testdata"

func Setup(t *testing.T, path string) *httptest.Server {
	t.Helper()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string][]byte)
	for _, info := range files {
		b, err := ioutil.ReadFile(filepath.Join(path, info.Name()))
		if err != nil {
			t.Fatal(err)
		}
		data[info.Name()] = b
	}
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		for k, b := range data {
			lower := strings.ToLower(req.URL.Path)
			if strings.Contains(lower, k) {
				res.Write(b)
				return
			}
		}
		http.Error(res, req.URL.Path, 500)
	}))
}

func TestOddHandler(t *testing.T) {
	ts := Setup(t, testdataPath)
	defer ts.Close()
	DefaultURL = fmt.Sprintf("%s/%s", ts.URL, "%s")
	o, err := New(&model.Odd{})
	if err != nil {
		t.Fatal(err)
	}
	handler := o.ServeHTTP
	req := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	for i := 0; i < 3; i++ {
		handler(w, req)
	}
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%s", body)
	}
}
