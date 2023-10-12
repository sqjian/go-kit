package jsonl_test

import (
	_ "embed"
	"github.com/sqjian/go-kit/encoding/jsonl"
	"testing"
)

//go:embed testdata/person.compressed.jsonl
var personCompressed []byte

//go:embed testdata/person.formatted.jsonl
var personFormatted []byte

type Person struct {
	Name string
	Age  int64
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
