package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Setup(t *testing.T, payload []byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(payload)
	}))
}

type Payload struct {
	Name  string
	Value int
}

func TestGetJSON(t *testing.T) {
	cases := map[string]struct {
		input  []byte
		expect Payload
	}{
		"base": {
			[]byte(`{"name":"test","value":1}`),
			Payload{Name: "test", Value: 1},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ts := Setup(t, tc.input)
			defer ts.Close()
			var out Payload
			if err := GetJSON(ts.URL, &out); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(out, tc.expect) {
				t.Errorf("want %#v got %#v", tc.expect, out)
			}
		})
	}
}
