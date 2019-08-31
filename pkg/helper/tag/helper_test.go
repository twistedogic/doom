package tag

import (
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
				name:  "Name",
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
