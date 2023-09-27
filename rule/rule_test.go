package main

import (
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
				//code: `get(fromJSON(ava_slot),"age")`,
				//code: `get({"name": "John", "age": 30}, "name") `,
				//code: `get({"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}, "age") `,
				//code: "keys(fromJSON(ava_slot))",
				//code: "get(keys(fromJSON(ava_slot)),0)",
				//code: "keys(fromJSON(ava_slot))|get(0)",
				//code: `fromJSON(ava_slot)|age`,
				code: `(get(fromJSON(ava_slot),"age")=="青年")?"年轻人":"老头子"`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    nil,
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
