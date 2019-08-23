package tap

import (
	"encoding/json"
	"testing"
)

type TestObj struct {
	Name   string
	Value  int
	Nested struct {
		Other string
	}
}

func TestNewSchema(t *testing.T) {
	var m Message
	if err := NewSchema(&TestObj{}, &m); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
