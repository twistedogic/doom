package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/structs"
	"github.com/google/go-cmp/cmp"
)

const testdataPath = "../../../../testdata"

func readFile(t *testing.T, name string) *os.File {
	t.Helper()
	filename := filepath.Join(testdataPath, name)
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestEntry(t *testing.T) {
	cases := map[string]struct {
		header   string
		value    string
		want     Entry
		json     string
		csv      map[string]string
		hasError bool
	}{
		"base": {
			header: "key,value",
			value:  "1,2",
			want:   Entry{[]string{"key", "value"}, []string{"1", "2"}},
			json:   `{"key":"1","value":"2"}`,
			csv: map[string]string{
				"key":   "1",
				"value": "2",
			},
			hasError: false,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := New(tc.header, tc.value)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.csv, got.MarshalCSV()); diff != "" {
				t.Fatal(diff)
			}
			b, err := json.Marshal(got)
			if (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.json, string(b)); diff != "" {
				t.Fatal(err)
			}
		})
	}
}

func TestDecoder(t *testing.T) {
	r := readFile(t, "csv")
	defer r.Close()
	ch := make(chan Entry)
	go func() {
		decoder := NewDecoder(r)
		if err := decoder.Decode(ch); err != nil {
			t.Fatal(err)
		}
	}()
	for entry := range ch {
		if structs.IsZero(entry) {
			t.Fatal(entry)
		}
	}
}
