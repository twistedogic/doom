package jockey

import (
	"testing"
)

func Test_Client_generateURL(t *testing.T) {
	cases := map[string]struct {
		base, bettype, want string
	}{
		"base": {
			Base,
			"fha",
			"https://bet.hkjc.com/football/getJSON.aspx?jsontype=odds_fha.aspx",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			client := New(tc.base, tc.bettype, -1)
			got := client.generateURL()
			if tc.want != got {
				t.Fatalf("want %s, got %s", tc.want, got)
			}
		})
	}
}
