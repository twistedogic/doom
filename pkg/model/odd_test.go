package model

import (
	"github.com/twistedogic/jsonpath"
	"testing"
)

func TestParseOdd(t *testing.T) {
	cases := []string{"fha", "hha", "had", "hft"}
	for i := range cases {
		name := cases[i]
		t.Run(name, func(t *testing.T) {
			data := ReadTestdata(t, name, "$.[*].matches.[*]")
			for _, b := range data {
				var o Odds
				if err := jsonpath.Unmarshal(b, &o); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
