package embedtpl

import (
	"embed"
	_ "embed"
	"testing"
)

//go:embed testdata/hello.txt
var f embed.FS

func TestFs(t *testing.T) {
	data, _ := f.ReadFile("testdata/hello.txt")
	t.Log(string(data))
}
