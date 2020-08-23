package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup(t *testing.T, body []byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(body)
	}))
}

func Test_request(t *testing.T) {
	cases := map[string]struct {
		body               []byte
		method             string
		isCancel, hasError bool
	}{
		"base": {
			body:     []byte(`something`),
			method:   "GET",
			isCancel: false,
			hasError: false,
		},
		"cancel": {
			body:     []byte(`something`),
			method:   "GET",
			isCancel: true,
			hasError: true,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ts := setup(t, tc.body)
			defer ts.Close()
			buf := new(bytes.Buffer)
			client := New(0)
			ctx, cancel := context.WithCancel(context.TODO())
			if tc.isCancel {
				cancel()
			}
			if err := client.Request(ctx, tc.method, ts.URL, nil, buf); (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			got := buf.Bytes()
			if !tc.hasError {
				if !bytes.Equal(tc.body, got) {
					t.Fatalf("want: %s, got: %s", string(tc.body), string(got))
				}
			}
		})
	}
}
