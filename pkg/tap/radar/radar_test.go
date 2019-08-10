package radar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

const testdataPath = "../../../testdata"

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

func TestClient(t *testing.T) {
	t.Skip()
	client := New(RadarURL)
	var out interface{}
	if err := client.GetMatchFullFeed(0, &out); err != nil {
		t.Fatal(err)
	} else {
		Write(t, filepath.Join(testdataPath, "fullfeed"), out)
	}
	if err := client.GetMatchDetail(18088911, &out); err != nil {
		t.Fatal(err)
	} else {
		Write(t, filepath.Join(testdataPath, "details"), out)
	}
	if err := client.GetBet(0, &out); err != nil {
		t.Fatal(err)
	} else {
		Write(t, filepath.Join(testdataPath, "bet"), out)
	}
	if err := client.GetLastMatches(54785, &out); err != nil {
		t.Fatal(err)
	} else {
		Write(t, filepath.Join(testdataPath, "team"), out)
	}
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
}
