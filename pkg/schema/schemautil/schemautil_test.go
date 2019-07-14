package schemautil

import (
	"reflect"
	"testing"
)

type Embedded struct {
	Name string
}

type Nested struct {
	Embedded Embedded
}

type Test struct {
	Name   string
	Value  int
	Nested Nested
}

func TestSetField(t *testing.T) {
	type args struct {
		obj   interface{}
		name  string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"base",
			args{
				obj:   &Test{},
				name:  "name",
				value: "something",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetField(tt.args.obj, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"base",
			args{Test{Name: "string", Value: 1, Nested: Nested{Embedded: Embedded{Name: "hidden"}}}},
			[]interface{}{
				Test{Name: "string", Value: 1, Nested: Nested{Embedded: Embedded{Name: "hidden"}}},
				Nested{Embedded: Embedded{Name: "hidden"}},
				Embedded{Name: "hidden"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unpack(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatToInt(t *testing.T) {
	input := 1.0
	out := floatToInt(reflect.ValueOf(input))
	t.Log(out)
}
