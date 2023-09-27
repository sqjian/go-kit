package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExecRule(t *testing.T) {
	type args struct {
		code string
		env  any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				code: `sprintf(greet, names[0])`,
				env: map[string]interface{}{
					"greet":   "Hello, %v!",
					"names":   []string{"world", "you"},
					"sprintf": fmt.Sprintf,
				},
			},
			want:    `Hello, world!`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecRule(tt.args.code, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}
