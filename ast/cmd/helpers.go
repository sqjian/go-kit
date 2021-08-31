package main

import (
	"bytes"
	"go/format"
	"go/token"
)

func node2String(node interface{}) (string, error) {
	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
