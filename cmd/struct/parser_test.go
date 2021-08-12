package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	var testData = `package geek

type geek1 struct {
	tag11 string
	tag21 string
}
type geek2 struct {
	tag12 string
	tag22 string
}
type geek3 struct {
	tag13 string
	tag23 string
}
`
	astParser := newAstParser()
	err := astParser.parseFile("", []byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(astParser.struct2string("geek2"))
}
