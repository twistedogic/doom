package radar

import (
	"encoding/json"
	"io/ioutil"
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
