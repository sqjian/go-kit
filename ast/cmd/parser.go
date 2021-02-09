package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sync"
)

type structType struct {
	name string
	node *ast.StructType
}

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
func (s *astParser) struct2string(structName string) (string, error) {
	expectedStruct, expectedStructErr := s.getStruct(structName)
	if expectedStructErr != nil {
		return "", expectedStructErr
	}
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "type %v ", expectedStruct.name)
	err := format.Node(&buf, s.fSet, expectedStruct.node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
func (s *astParser) getStruct(structName string) (*structType, error) {
	s.rwMu.RLock()
	defer s.rwMu.RUnlock()

	var expectedStruct *structType
	for _, st := range s.structs {
		if st.name == structName {
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
