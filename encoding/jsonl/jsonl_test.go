package jsonl_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/sqjian/go-kit/encoding/jsonl"
	"testing"
)

//go:embed testdata/person.compressed.jsonl
var personCompressed []byte

//go:embed testdata/person.formatted.jsonl
var personFormatted []byte

//go:embed testdata/person.formatted.commented.jsonl
var personFormattedCommented []byte

//go:embed testdata/person.formatted.commented.extra.jsonl
var personFormattedCommentedExtra []byte

//go:embed testdata/dev.jsonl
var dev []byte

type Person struct {
	Name string
	Age  int64
}

type Dev struct {
	Input  string `json:"input" `
	Target string `json:"target"`
}

func TestDecodeDev(t *testing.T) {
	type args struct {
		data       []byte
		ptrToSlice *Dev
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "dev",
			args: args{
				data:       dev,
				ptrToSlice: &Dev{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonl.Decode(bytes.NewReader(tt.args.data), func(jsonBuffer []byte) error {
				return json.Unmarshal(jsonBuffer, tt.args.ptrToSlice)
			})
			if err != nil {
				t.Fatalf("Unmarshal returns error:%v,data:%s", err, tt.args.data)
			}
		})
	}
}

func TestUnmarshalDev(t *testing.T) {
	type args struct {
		data       []byte
		ptrToSlice *[]Dev
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "dev",
			args: args{
				data:       dev,
				ptrToSlice: &[]Dev{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonl.Unmarshal(tt.args.data, tt.args.ptrToSlice)
			if err != nil {
				t.Fatalf("Unmarshal returns error:%v,data:%s", err, tt.args.data)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data       []byte
		ptrToSlice *[]Person
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "personCompressed",
			args: args{
				data:       personCompressed,
				ptrToSlice: &[]Person{},
			},
		},
		{
			name: "personFormatted",
			args: args{
				data:       personFormatted,
				ptrToSlice: &[]Person{},
			},
		},
		{
			name: "personFormattedCommented",
			args: args{
				data:       personFormattedCommented,
				ptrToSlice: &[]Person{},
			},
		},
		{
			name: "personFormattedCommentedExtra",
			args: args{
				data:       personFormattedCommentedExtra,
				ptrToSlice: &[]Person{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonl.Unmarshal(tt.args.data, tt.args.ptrToSlice)

			if err != nil {
				t.Fatalf("Unmarshal returns error:%v,data:%s", err, tt.args.data)
			}
			if len(*tt.args.ptrToSlice) != 3 {
				t.Fatalf("Expected 2 objects in slice, got %v", len(*tt.args.ptrToSlice))
			}
			if (*tt.args.ptrToSlice)[0].Name != "Paul" {
				t.Fatalf("Unexpected value in first object Name field: %v", (*tt.args.ptrToSlice)[0].Name)
			}
			if (*tt.args.ptrToSlice)[1].Name != "John" {
				t.Fatalf("Unexpected value in first object Name field: %v", (*tt.args.ptrToSlice)[0].Name)
			}
			if (*tt.args.ptrToSlice)[2].Name != "Maria" {
				t.Fatalf("Unexpected value in first object Name field: %v", (*tt.args.ptrToSlice)[0].Name)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	type args struct {
		ptrToSlice any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "person",
			args: args{
				ptrToSlice: []*Person{
					{Name: "Paul", Age: 20},
					{Name: "John", Age: 30},
				},
			},
		},
		{
			name: "personPtr",
			args: args{
				ptrToSlice: &[]*Person{
					{Name: "Paul", Age: 20},
					{Name: "John", Age: 30},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := jsonl.Marshal(tt.args.ptrToSlice)
			if err != nil {
				t.Fatal("Marshal returns error: ", err)
			}
			t.Log(string(data))
		})
	}
}
