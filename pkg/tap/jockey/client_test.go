package jockey

import (
	"testing"
)

const testdataPath = "../../../testdata"

func Test_toQueryString(t *testing.T) {
	cases := map[string]struct {
		kv   map[string]string
		want string
	}{
		"base": {
			kv: map[string]string{
				"key":   "value",
				"other": "thing",
			},
			want: "key=value&other=thing",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := toQueryString(tc.kv)
			if tc.want != got {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
