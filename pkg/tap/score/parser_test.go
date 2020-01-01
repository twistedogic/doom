package score

import (
	"io/ioutil"
	"testing"

	"github.com/twistedogic/doom/testutil"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		filename string
		repeat   int
		want     []Score
	}{
		"base": {
			filename: "stream",
			repeat:   2,
			want:     []Score{},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			f := testutil.Open(t, testdataPath, tc.filename, tc.repeat)
			b, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			got, err := Parse(b)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(got)
		})
	}
}
