package fs

import (
	_ "embed"
	"testing"
)

//go:embed testdata/hello.txt
var b []byte

func TestByte(t *testing.T) {
	t.Log(string(b))
}
