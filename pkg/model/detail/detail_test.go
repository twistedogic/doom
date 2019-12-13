package detail

import (
	"testing"

	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func TestTransform(t *testing.T) {
	cases := map[string]struct {
		name string
	}{
		"base": {"details"},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			f := testutil.Open(t, testdataPath, tc.name, 2)
			target := testutil.NewMockTarget(t, false, false)
			if err := Transform(f, target); err != nil {
				t.Fatal(err)
			}
			if len(target.Bytes()) == 0 {
				t.Fatal(string(target.Bytes()))
			}
			t.Log(string(target.Bytes()))
		})
	}
}
