package main

import (
	"fmt"
	"go/ast"
)

type structType struct {
	name string
	node *ast.StructType
}

func (s *structType) getFiled() ([]string, error) {
	for i := 0; i < s.node.Fields.NumFields(); i++ {
		fmt.Printf("filed:%+v\n", fmt.Sprint(node2String(s.node.Fields.List[i].Names[0])))
	}
	return nil, nil
}
