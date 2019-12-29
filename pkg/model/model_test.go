package model

import (
	"bufio"
	"context"
	"errors"
	"io"
	"testing"
	"time"

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

func setupMockTransformer(t *testing.T, ttype Type, hasError bool) Transformer {
	t.Helper()
	transform := func(r io.Reader, e Encoder) error {
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
			input: []byte("line1\nline2\nline3\n"),
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
			input: []byte("line1\nline2\nline3\n"),
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
			transformers := make([]Transformer, 0, len(tc.transformers))
			for _, tt := range tc.transformers {
				transformers = append(transformers, setupMockTransformer(
					t, tt.ttype, tt.hasError,
				))
			}
			m := New(s)
			ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
			defer cancel()
			if _, err := m.Write(tc.input); err != nil {
				t.Fatal(err)
			}
			if err := m.Update(ctx, transformers...); (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			if diff := cmp.Diff(s.Content(), tc.want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
