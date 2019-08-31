package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigNew(t *testing.T) {
	cases := map[string]struct {
		input    []byte
		want     Config
		hasError bool
	}{
		"base": {
			input: []byte(`tasks:
- tap:
    name: tap 
    config:
      key: a
      value: 2
  target:
    name: target
    config:
      key: b
      value: 3
  schedule: '* * * * *'`),
			want: Config{
				Tasks: []Task{
					{
						Tap:      Setting{Name: "tap", Config: map[string]string{"key": "a", "value": "2"}},
						Target:   Setting{Name: "target", Config: map[string]string{"key": "b", "value": "3"}},
						Schedule: "* * * * *",
					},
				},
			},
			hasError: false,
		},
		"missing field": {
			input: []byte(`tasks:
- tap:
    name: tap 
    config:
      key: a
      value: 2`),
			want: Config{
				Tasks: []Task{
					{
						Tap: Setting{Name: "tap", Config: map[string]string{"key": "a", "value": "2"}},
					},
				},
			},
			hasError: false,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got, err := New(tc.input)
			if (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" && err == nil {
				t.Fatal(diff)
			}
		})
	}
}
