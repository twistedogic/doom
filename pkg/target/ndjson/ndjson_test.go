package ndjson

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Setup(t *testing.T, name string) string {
	t.Helper()
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestTarget(t *testing.T) {
	type base struct {
		Name  string
		Value int
	}
	cases := map[string]struct {
		input []base
		want  string
	}{
		"base": {
			input: []base{
				{"test", 1},
				{"test", 1},
				{"test", 1},
			},
			want: `{"Name":"test","Value":1}
{"Name":"test","Value":1}
{"Name":"test","Value":1}
`,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			filename := Setup(t, name)
			defer os.Remove(filename)
			target, err := New(filename)
			if err != nil {
				t.Fatal(err)
			}
			if err := target.BulkUpsert(tc.input); err != nil {
				t.Fatal(err)
			}
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, string(b)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
