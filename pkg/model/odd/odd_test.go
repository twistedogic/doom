package odd

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func TestTransform(t *testing.T) {
	cases := map[string]struct {
		filename string
		repeat   int
	}{
		"fha": {
			filename: "fha",
			repeat:   2,
		},
		"had": {
			filename: "had",
			repeat:   2,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			f := testutil.Open(t, testdataPath, tc.filename, tc.repeat)
			s := testutil.NewMockStore(t, make(map[string][]byte), false)
			m := model.New(s)
			ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
			defer cancel()
			if _, err := io.Copy(m, f); err != nil {
				t.Fatal(err)
			}
			if err := m.Update(ctx, Transform); err != nil {
				t.Fatal(err)
			}
			for k, v := range s.Content() {
				if !strings.HasPrefix(k, string(Type)) {
					t.Fatal(k)
				}
				var model Model
				if err := json.Unmarshal(v, &model); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
