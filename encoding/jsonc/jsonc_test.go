package jsonc

import (
	_ "embed"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/helper"
	"strings"
	"testing"
)

//go:embed testdata/person.jsonc
var personCommented []byte

//go:embed testdata/comments.go
var comments []byte

//go:embed testdata/person.formatted.commented.jsonl
var personFormattedCommented []byte

func Test_Comments(t *testing.T) {
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
				data: comments,
			},
		},
		{
			name: "test2",
			args: args{
				data: personFormattedCommented,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimCommentWrapper(helper.BytesToRunes(tt.args.data))
			spew.Dump(string(got))
			if strings.Contains(string(got), "//") {
				t.Fatalf("still containing comment chars")
			}
		})
	}
}

func Test_TrimCommentWrapper(t *testing.T) {
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
			got := TrimCommentWrapper(helper.BytesToRunes(tt.args.data))
			v := make(map[string]any)
			if err := json.Unmarshal(helper.RunesToBytes(got), &v); err != nil {
				t.Fatalf("unmarshal failed,err:%v", err)
			}
			if (v["Name"].(string)) != "Maria" {
				t.Fatalf("Unexpected value in object Name field: %v", v["Name"])
			}
			if (v["Age"].(float64)) != 30 {
				t.Fatalf("Unexpected value in object Name field: %v", v["Age"])
			}
			if strings.Contains(string(got), "//") {
				t.Fatalf("still containing comment chars")
			}
		})
	}
}
