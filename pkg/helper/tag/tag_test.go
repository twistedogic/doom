package tag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type obj struct {
	Id int `odd:"id"`
}

func TestGetTaggedFields(t *testing.T) {
	OddTag := "odd"
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"No Tag",
			args{
				struct {
					Name string
				}{
					"hi",
				},
			},
			[]string{},
		}, {
			"flat",
			args{
				obj{11},
			},
			[]string{"Id"},
		}, {
			"ignore nested",
			args{
				struct {
					Obj obj
				}{
					obj{12},
				},
			},
			[]string{},
		}, {
			"tagged struct",
			args{
				struct {
					Obj obj `odd:"obj"`
				}{
					obj{12},
				},
			},
			[]string{"Obj"},
		}, {
			"parse empty",
			args{
				struct {
					Test int `odd:"obj"`
				}{},
			},
			[]string{"Test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTaggedFields(tt.args.s, OddTag)
			names := make([]string, len(got))
			for i, v := range got {
				names[i] = v.Name
			}
			if !cmp.Equal(names, tt.want) {
				t.Errorf("GetTaggedFields() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
