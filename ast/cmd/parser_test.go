package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	astParser := newAstParser()
	err := astParser.parseFile("test_data/geek.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(astParser.struct2string("geek2"))
}
