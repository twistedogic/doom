package csv

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshal(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		type base struct {
			Name  string
			Value int
		}
		input := []base{
			{"test", 1},
			{"test", 1},
			{"test", 1},
		}
		expect := "name,value\ntest,1\ntest,1\ntest,1\n"
		buf := &bytes.Buffer{}
		w := csv.NewWriter(buf)
		if err := Marshal(w, input, true); err != nil {
			t.Fatal(err)
		}
		output := buf.String()
		if diff := cmp.Diff(expect, output); diff != "" {
			t.Fatal(diff)
		}
	})
}
