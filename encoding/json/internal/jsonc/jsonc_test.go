package jsonc

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/person.jsonc
var personCommented []byte

func Test_translate(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				data: personCommented,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Translate(tt.args.data)
			v := make(map[string]any)
			if err := json.Unmarshal(got, &v); err != nil {
				t.Fatalf("unmarshal failed,err:%v", err)
			}
			if (v["Name"].(string)) != "Maria" {
				t.Fatalf("Unexpected value in object Name field: %v", v["Name"])
			}
			if (v["Age"].(float64)) != 30 {
				t.Fatalf("Unexpected value in object Name field: %v", v["Age"])
			}
		})
	}
}
