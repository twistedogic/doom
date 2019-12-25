package model

import (
	"bufio"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/twistedogic/doom/testutil"
)

type mockData struct {
	Type Type
	Data []byte
}

func (m mockData) Item(i *Item) error {
	i.Key, i.Type, i.Data = string(m.Data), m.Type, m.Data
	return nil
}

func setupMockTransform(t *testing.T, ttype Type, hasError bool) func(io.Reader, Encoder) error {
	t.Helper()
	return func(r io.Reader, e Encoder) error {
		if hasError {
			return errors.New("error")
		}
		buf := bufio.NewReader(r)
		for {
			line, _, err := buf.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := e.Encode(mockData{ttype, line}); err != nil {
				return err
			}
		}
		return nil
	}
}

func TestModeler(t *testing.T) {
	cases := map[string]struct {
		input        []byte
		transformers []struct {
			ttype    Type
			hasError bool
		}
		hasError bool
	}{
		"base": {
			input: []byte("line1\nline2\nline3"),
			transformers: []struct {
				ttype    Type
				hasError bool
			}{
				{Type("test1"), false},
				{Type("test2"), false},
			},
			hasError: false,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			s := testutil.NewMockStore(t, make(map[string][]byte), false)
			transformers := make([]Transformer, len(tc.transformers))
			for i, tt := range tc.transformers {
				transformers[i] = setupMockTransform(t, tt.ttype, tt.hasError)
			}
			m := New(s)
			ctx := context.TODO()
			if _, err := m.Write(tc.input); err != nil {
				t.Fatal(err)
			}
			if err := m.Update(ctx, transformers...); (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			t.Log(s.Content())
		})
	}
}
