package jc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

const testdataPath = "../../../testdata"

func Write(t *testing.T, filePath string, value interface{}) {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filePath, b, 0644); err != nil {
		t.Fatal(err)
	}
}

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

func TestClient(t *testing.T) {
	t.Skip()
	client := New(JcURL)
	cases := []string{"fha", "hha", "had", "hft"}
	for i := range cases {
		name := cases[i]
		t.Run(name, func(t *testing.T) {
			var out interface{}
			if err := client.GetOdd(name, &out); err != nil {
				t.Fatal(err)
			} else {
				Write(t, filepath.Join(testdataPath, name), out)
			}
		})
	}
}

func TestGetOdds(t *testing.T) {
	ts := Setup(t, "fha")
	defer ts.Close()
	client := New(fmt.Sprintf("%s/%s", ts.URL, "%s"))
	results, err := client.GetOdds()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fail()
	}
}
