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
