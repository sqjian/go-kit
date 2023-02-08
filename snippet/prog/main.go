package main

import (
	"fmt"
	"github.com/sqjian/go-kit/snippet/prog/cmd"
	"github.com/sqjian/go-kit/splash"
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
