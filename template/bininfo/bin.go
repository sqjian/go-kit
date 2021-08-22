package bininfo

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	GitTag         = "Unknown"
	GitCommitLog   = "Unknown"
	GitStatus      = "Unknown"
	BuildTime      = "Unknown"
	BuildGoVersion = "Unknown"
)

func StringifySingleLine() string {
	return fmt.Sprintf("GitTag=%s. GitCommitLog=%s. GitStatus=%s. BuildTime=%s. GoVersion=%s. runtime=%s/%s.",
		GitTag, GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func StringifyMultiLine() string {
	return fmt.Sprintf("GitTag=%s\nGitCommitLog=%s\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		GitTag, GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func beauty() {
	if GitStatus == "" {
		GitStatus = "cleanly"
	} else {
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}
