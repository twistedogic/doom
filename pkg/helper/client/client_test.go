package client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/pkg/config"
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

func TestExtractJsonPath(t *testing.T) {
	path := "$.name"
	input := map[string]string{
		"name": "test",
	}
	want := [][]byte{[]byte(`"test"`)}
	got, err := ExtractJsonPath(input, path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, want) {
		for _, item := range got {
			t.Logf("%s", item)
		}
		t.Errorf("want %#v got %#v", want, got)
	}
}

func TestLoad(t *testing.T) {
	got := &Client{}
	input := config.Setting{
		Name: "test",
		Config: map[string]string{
			"baseurl": "test",
			"rate":    "10",
		},
	}
	want := New("test", 10)
	if err := got.Load(input); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %#v, got: %#v", want, got)
	}
}
