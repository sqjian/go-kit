package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sync"
)

type astParser struct {
	rwMu    sync.RWMutex
	fSet    *token.FileSet
	structs map[token.Pos]*structType
}

func newAstParser() *astParser {
	return &astParser{
		rwMu: sync.RWMutex{},
		structs: make(
			map[token.Pos]*structType, 0,
		),
		fSet: token.NewFileSet(),
	}
}

func (s *astParser) extractStruct(name string) (*structType, error) {
	s.rwMu.RLock()
	defer s.rwMu.RUnlock()

	var expectedStruct *structType
	for _, st := range s.structs {
		if st.name == name {
			expectedStruct = st
		}
	}
	if expectedStruct == nil {
		return nil, fmt.Errorf("struct name does not exist")
	}
	return expectedStruct, nil
}
func (s *astParser) parseFile(filename string, src interface{}) error {
	s.rwMu.Lock()
	defer s.rwMu.Unlock()

	node, err := parser.ParseFile(
		s.fSet,
		filename,
		src,
		parser.ParseComments,
	)
	if err != nil {
		return err
	}
	structs := make(map[token.Pos]*structType, 0)
	collectStructs := func(n ast.Node) bool {
		var t ast.Expr
		var structName string
		switch x := n.(type) {
		case *ast.TypeSpec:
			{
				if x.Type == nil {
					return true
				}
				structName = x.Name.Name
				t = x.Type
			}
		}
		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}
		structs[x.Pos()] = &structType{
			name: structName,
			node: x,
		}
		return true
	}
	ast.Inspect(node, collectStructs)
	s.structs = structs
	return nil
}
