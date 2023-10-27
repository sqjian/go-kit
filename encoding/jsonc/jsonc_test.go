package jsonc

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"sync"
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
			got := TrimCommentWrapper(tt.args.data)
			spew.Dump(string(got))
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
			got := TrimCommentWrapper(tt.args.data)
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

func Test_trimComment(t *testing.T) {
	var wg sync.WaitGroup

	unProcessed := func() chan byte {
		ch := make(chan byte)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, i := range personCommented {
				ch <- i
			}
			close(ch)
		}()
		return ch
	}()

	processed := make(chan byte)

	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = TrimComment(context.Background(), unProcessed, processed)
		}()
	}

	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			var rst []byte
			for ch := range processed {
				rst = append(rst, ch)
			}
			spew.Dump(string(rst))
		}()
	}

	wg.Wait()
}
