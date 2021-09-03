package splash

import (
	"bytes"
	"fmt"
	"html/template"
)

const splash = `
--------------------------------
   / \__
  (    @\___
  /         O
 /   (_____/
/_____/   U

GitTag: {{if .GitTag}}{{.GitTag}}{{else}}Unknown{{end}}
GitCommitLog: {{if .GitCommitLog}}{{.GitCommitLog}}{{else}}Unknown{{end}}
GitStatus: {{if .GitStatus}}{{.GitStatus}}{{else}}Unknown{{end}}
BuildTime: {{if .BuildTime}}{{.BuildTime}}{{else}}Unknown{{end}}
BuildGoVersion: {{if .BuildGoVersion}}{{.BuildGoVersion}}{{else}}Unknown{{end}}
--------------------------------
`

var (
	BinInfo Tag
	tpl     = template.Must(template.New("splash").Parse(splash))
)

type Tag struct {
	GitTag         string
	GitCommitLog   string
	GitStatus      string
	BuildTime      string
	BuildGoVersion string
}

func Stringify() string {
	fmt.Println("BinInfo:", BinInfo)

	buf := bytes.NewBuffer(nil)

	err := tpl.Execute(buf, BinInfo)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
