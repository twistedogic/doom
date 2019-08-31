package flatten

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlattenKey(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		type base struct {
			Name  string
			Value int
		}
		input := base{"test", 1}
		expect := []map[string]string{
			{"name": "test", "value": "1"},
		}
		output := FlattenKey(input, ".")
		if diff := cmp.Diff(expect, output); diff != "" {
			t.Fatal(diff)
		}
	})
	t.Run("slice", func(t *testing.T) {
		type base struct {
			Name  string
			Value int
		}
		input := []base{{"test", 1}}
		expect := []map[string]string{
			{"name": "test", "value": "1"},
		}
		output := FlattenKey(input, ".")
		if diff := cmp.Diff(expect, output); diff != "" {
			t.Fatal(diff)
		}
	})
	t.Run("nested", func(t *testing.T) {
		type nested struct {
			Value float64
		}
		type base struct {
			Name      string
			Value     int
			Something []string
			Nested    nested
		}
		input := []base{{"test", 1, []string{"test", "other", "this thing"}, nested{1.01}}}
		expect := []map[string]string{
			{"name": "test", "value": "1", "something": "test,other,this thing", "nested.value": "1.01"},
		}
		output := FlattenKey(input, ".")
		if diff := cmp.Diff(expect, output); diff != "" {
			t.Fatal(diff)
		}
	})
}
