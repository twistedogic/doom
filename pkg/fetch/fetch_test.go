package fetch

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

const testdataPath = "../../testdata"

func Setup(t *testing.T, file string) *httptest.Server {
	t.Helper()
	b, err := ioutil.ReadFile(filepath.Join(testdataPath, file))
	if err != nil {
		t.Fatal(err)
	}
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(b)
	}))
}

func TestGetMatch(t *testing.T) {
	ts := Setup(t, "fullfeed")
	defer ts.Close()
	f := New(ts.URL)
	results, err := f.GetMatch(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fail()
	}
	t.Logf("%#v", results)
}

func TestGetDetail(t *testing.T) {
	ts := Setup(t, "details")
	defer ts.Close()
	f := New(ts.URL)
	results, err := f.GetDetail(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fail()
	}
	t.Logf("%#v", results)
}
