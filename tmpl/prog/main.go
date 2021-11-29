package main

import (
	"fmt"
	"github.com/sqjian/go-kit/splash"
	"github.com/sqjian/go-kit/tmpl/prog/cmd"
)

var (
	GitTag         string
	BuildTime      string
	GitCommitLog   string
	BuildGoVersion string
)

func init() {
	fmt.Println(splash.Stringify(GitTag, GitCommitLog, BuildTime, BuildGoVersion))
}
func main() {
	cmd.Execute()
}
