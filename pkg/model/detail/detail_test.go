package detail

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/fatih/structs"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func TestTransform(t *testing.T) {
	cases := map[string]struct {
		filename string
	}{
		"base": {
			filename: "details",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			log.SetOutput(buf)
			f := testutil.Open(t, testdataPath, tc.filename, 1)
			s := testutil.NewMockStore(t, make(map[string][]byte), false)
			m := model.New(s, Transform)
			if _, err := io.Copy(m, f); err != nil {
				t.Fatal(err)
			}
			if len(buf.Bytes()) != 0 {
				t.Fatal(buf.String())
			}
			for k, v := range s.Content() {
				if !strings.HasPrefix(k, string(Type)) {
					t.Fatal(k)
				}
				var model Model
				if err := json.Unmarshal(v, &model); err != nil {
					t.Fatal(err)
				}
				if structs.HasZero(model) {
					t.Fatal(model)
				}
			}
		})
	}
}
