package configor

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestLoadJsonContents(t *testing.T) {
	type args struct {
		obj  any
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				obj: &struct {
					Name   string `validate:"required"`
					Age    int    `validate:"gte=10,lte=130"`
					Gender int
				}{},
				data: []byte(`{"name": "sqjian","age": 12}`),
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadJsonContents(tt.args.obj, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("LoadJsonContents() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				spew.Dump(tt.args.obj)
			}
		})
	}
}
