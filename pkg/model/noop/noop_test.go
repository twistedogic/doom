package noop

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/testutil"
)

func TestTransform(t *testing.T) {
	cases := map[string]struct {
		content []byte
	}{
		"base": {
			content: []byte(`test`),
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			log.SetOutput(buf)
			s := testutil.NewMockStore(t, make(map[string][]byte), false)
			m := model.New(s, Transform)
			if err := m.Write(tc.content); err != nil {
				t.Fatal(err)
			}
			if len(buf.Bytes()) != 0 {
				t.Fatal(buf.String())
			}
			for k, v := range s.Content() {
				if !strings.HasPrefix(k, string(Type)) {
					t.Fatal(k)
				}
				if string(v) != string(tc.content) {
					t.Fatal(v)
				}
				t.Logf("%s, %s", k, string(v))
			}
		})
	}
}
