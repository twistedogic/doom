package model

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/fatih/structs"
	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/jsonpath"
)

const testdataPath = "../../testdata"

func ReadTestFile(t *testing.T, file string) []byte {
	t.Helper()
	b, err := ioutil.ReadFile(filepath.Join(testdataPath, file))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func ReadTestdata(t *testing.T, file, pattern string) [][]byte {
	t.Helper()
	b := ReadTestFile(t, file)
	var container interface{}
	if err := json.Unmarshal(b, &container); err != nil {
		t.Fatal(err)
	}
	value, err := jsonpath.Lookup(pattern, container)
	if err != nil {
		t.Fatal(err)
	}
	values := helper.FlattenDeep(value)
	out := make([][]byte, len(values))
	for i, v := range values {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		out[i] = b
	}
	return out
}

func TestParseMatch(t *testing.T) {
	data := ReadTestdata(t, "fullfeed", "$.doc[*].data[*].realcategories[*].tournaments[*].matches")
	for _, b := range data {
		var m Match
		if err := jsonpath.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		if structs.IsZero(m) {
			t.Fatalf("not parsed %#v", m)
		}
	}
}

func TestParseDetail(t *testing.T) {
	data := ReadTestdata(t, "details", "$.doc[*].data[*]")
	for _, b := range data {
		var d Detail
		if err := jsonpath.Unmarshal(b, &d); err != nil {
			t.Fatal(err)
		}
		if structs.IsZero(d) {
			t.Fatalf("not parsed %#v", d)
		}
	}
}

func TestParseOdd(t *testing.T) {
	cases := []string{"fha", "hha", "had", "hft"}
	for i := range cases {
		name := cases[i]
		t.Run(name, func(t *testing.T) {
			data := ReadTestdata(t, name, "$.[*].matches.[*]")
			for _, b := range data {
				var o Odds
				if err := jsonpath.Unmarshal(b, &o); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
