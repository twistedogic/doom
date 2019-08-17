package radar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testdataPath = "../../../testdata"

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
	client := New(RadarURL, -1)
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
	ts := Setup(t, testdataPath)
	defer ts.Close()
	f := New(ts.URL, -1)
	results, err := f.GetMatch(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fail()
	}
}

func TestGetDetail(t *testing.T) {
	ts := Setup(t, testdataPath)
	defer ts.Close()
	f := New(ts.URL, -1)
	results, err := f.GetDetail(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fail()
	}
}

type mockTarget struct {
	t *testing.T
}

func NewMockTarget(t *testing.T) mockTarget {
	t.Helper()
	return mockTarget{t}
}

func (m mockTarget) UpsertItem(interface{}) error { return nil }
func (m mockTarget) BulkUpsert(interface{}) error { return nil }
func (m mockTarget) GetLastUpdate() time.Time     { return time.Now() }

func TestUpdate(t *testing.T) {
	ts := Setup(t, testdataPath)
	defer ts.Close()
	target := NewMockTarget(t)
	f := New(ts.URL, -1)
	if err := f.Update(target); err != nil {
		t.Fatal(err)
	}
}
