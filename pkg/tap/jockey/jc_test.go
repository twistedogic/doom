package jockey

import (
	"context"
	"testing"

	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

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

func TestUpdate(t *testing.T) {
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	ctx := context.TODO()
	target := testutil.NewMockTarget(t, false, false)
	for i := range BetTypes {
		bet := BetTypes[i]
		t.Run(bet, func(t *testing.T) {
			f := New(ts.URL, bet, -1)
			if err := f.Update(ctx, target); err != nil {
				t.Fatal(err)
			}
		})
	}

}
