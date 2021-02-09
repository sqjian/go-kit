package main

import (
	"testing"
)

func TestStructLab(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	var inst *structType
	{
		astParser := newAstParser()
		err := astParser.parseFile("test_data/geek.go", nil)
		if err != nil {
			t.Fatal(err)
		}
		tmp, tmpErr := astParser.extractStruct("geek2")
		checkErr(tmpErr)
		inst = tmp
	}

	{
		t.Log(node2String(inst.node))
	}

	{
		_, _ = inst.getFiled()
	}
}
