package splash

import (
	"bytes"
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
BuildTime: {{if .BuildTime}}{{.BuildTime}}{{else}}Unknown{{end}}
BuildGoVersion: {{if .BuildGoVersion}}{{.BuildGoVersion}}{{else}}Unknown{{end}}
--------------------------------
`

var (
	tpl = template.Must(template.New("splash").Parse(splash))
)

type Tag struct {
	GitTag         string
	GitCommitLog   string
	BuildTime      string
	BuildGoVersion string
}

func Stringify(GitTag string, GitCommitLog string, BuildTime string, BuildGoVersion string) string {

	buf := bytes.NewBuffer(nil)

	err := tpl.Execute(buf, Tag{
		GitTag:         GitTag,
		GitCommitLog:   GitCommitLog,
		BuildTime:      BuildTime,
		BuildGoVersion: BuildGoVersion,
	})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
