package model

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func setupMockTransformFunc(t *testing.T, ttype Type, hasError bool) TransformFunc {
	t.Helper()
	transform := func(b []byte, e Encoder) error {
		if hasError {
			return errors.New("error")
		}
		for _, line := range bytes.Split(b, []byte("\n")) {
			data := mockData{ttype, line}
			if err := e.Encode(data); err != nil {
				return err
			}
		}
		return nil
	}
	return transform
}

func TestModeler(t *testing.T) {
	cases := map[string]struct {
		input        []byte
		transformers []struct {
			ttype    Type
			hasError bool
		}
		hasError bool
		want     map[string][]byte
	}{
		"single transformer": {
			input: []byte("line1\nline2\nline3"),
			transformers: []struct {
				ttype    Type
				hasError bool
			}{
				{Type("test1"), false},
			},
			hasError: false,
			want: map[string][]byte{
				"test1:line1": []byte("line1"),
				"test1:line2": []byte("line2"),
				"test1:line3": []byte("line3"),
			},
		},
		"multi-transformers": {
			input: []byte("line1\nline2\nline3"),
			transformers: []struct {
				ttype    Type
				hasError bool
			}{
				{Type("test1"), false},
				{Type("test2"), false},
			},
			hasError: false,
			want: map[string][]byte{
				"test1:line1": []byte("line1"),
				"test1:line2": []byte("line2"),
				"test1:line3": []byte("line3"),
				"test2:line1": []byte("line1"),
				"test2:line2": []byte("line2"),
				"test2:line3": []byte("line3"),
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			s := testutil.NewMockStore(t, make(map[string][]byte), false)
			transformers := make([]TransformFunc, 0, len(tc.transformers))
			for _, tt := range tc.transformers {
				transformers = append(transformers, setupMockTransformFunc(
					t, tt.ttype, tt.hasError,
				))
			}
			m := New(s, transformers...)
			if _, err := m.Write(tc.input); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(s.Content(), tc.want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
