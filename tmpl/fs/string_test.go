package fs

import (
	_ "embed"
	"testing"
)

//go:embed testdata/hello.txt
var s string

func TestString(t *testing.T) {
	t.Log(s)
}
