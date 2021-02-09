package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	astParser := newAstParser()
	err := astParser.parseFile("test_data/geek.go", nil)
	if err != nil {
		t.Fatal(err)
	}

	structInst, structInstErr := astParser.extractStruct("geek2")
	checkErr(structInstErr)

	t.Log(node2String(structInst.node))
}

func TestGetExporter(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	astParser := newAstParser()
	err := astParser.parseFile("test_data/geek.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	tmp, tmpErr := astParser.extractStruct("geek2")
	checkErr(tmpErr)

	{ // filed
		t.Logf("filed1:%+v", fmt.Sprint(node2String(tmp.node.Fields.List[0].Names[0])))
		t.Logf("filed2:%+v", fmt.Sprint(node2String(tmp.node.Fields.List[1].Names[0])))
	}
	{ // type
		t.Logf("type:%+v", fmt.Sprint(node2String(tmp.node.Fields.List[0].Type)))
	}
	{
		// comment
		t.Log(tmp.node.Fields.List[0].Comment.Text())
	}
}

func TestAstLab(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	astParser := newAstParser()
	err := astParser.parseFile("test_data/geek.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	tmp, tmpErr := astParser.extractStruct("geek2")
	checkErr(tmpErr)
	{ // filed
		t.Logf("filed1:%+v", fmt.Sprint(node2String(tmp.node.Fields.List[0].Names[0])))
		t.Logf("filed2:%+v", fmt.Sprint(node2String(tmp.node.Fields.List[1].Names[0])))
	}
}
