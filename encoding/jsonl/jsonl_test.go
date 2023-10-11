package jsonl_test

import (
	"bytes"
	_ "embed"
	"github.com/sqjian/go-kit/encoding/jsonl"
	"io"
	"strings"
	"testing"
)

//go:embed testdata/person.compressed.jsonl
var personCompressed string

//go:embed testdata/person.formatted.jsonl
var personFormatted string

type Person struct {
	Name string
	Age  int64
}

func TestDecode(t *testing.T) {
	type args struct {
		reader     io.Reader
		ptrToSlice *[]Person
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				reader:     strings.NewReader(personCompressed),
				ptrToSlice: &[]Person{},
			},
		},
		{
			name: "test2",
			args: args{
				reader:     strings.NewReader(personFormatted),
				ptrToSlice: &[]Person{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonl.Decode(tt.args.reader, tt.args.ptrToSlice)

			if err != nil {
				t.Fatal("Decode returns error: ", err)
			}
			if len(*tt.args.ptrToSlice) != 2 {
				t.Fatalf("Expected 2 objects in slice, got %v", len(*tt.args.ptrToSlice))
			}
			if (*tt.args.ptrToSlice)[0].Name != "Paul" {
				t.Fatalf("Unexpected value in first object Name field: %v", (*tt.args.ptrToSlice)[0].Name)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		ptrToSlice any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
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
			buf := &bytes.Buffer{}
			err := jsonl.Encode(buf, tt.args.ptrToSlice)
			if err != nil {
				t.Fatal("Encode returns error: ", err)
			}
			t.Log(buf.String())
		})
	}
}

func TestDecodeWrongTypes(t *testing.T) {
	type args struct {
		ptrToSlice any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				ptrToSlice: &map[Person]int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := jsonl.Encode(buf, tt.args.ptrToSlice)
			if err == nil {
				t.Fatal("Decode doesn't returns error")
			}
			if err.Error() != "expected pointer to slice, got pointer to map" {
				t.Fatalf("Decode return wrong error: %v", err)
			}
		})
	}
}
