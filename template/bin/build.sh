#set -x

GitTag=$(git tag --sort=version:refname | tail -n 1)
GitCommitLog=$(git log --pretty=oneline -n 1)
GitCommitLog=${GitCommitLog//\'/\"}
GitStatus=$(git status -s)
BuildTime=$(date +'%Y.%m.%d.%H%M%S')
BuildGoVersion=$(go version)

LDFlags=" \
    -X 'github.com/sqjian/go-kit/splash.BinInfo.GitTag=${GitTag}' \
    -X 'github.com/sqjian/go-kit/splash.BinInfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/sqjian/go-kit/splash.BinInfo.GitStatus=${GitStatus}' \
    -X 'github.com/sqjian/go-kit/splash.BinInfo.BuildTime=${BuildTime}' \
    -X 'github.com/sqjian/go-kit/splash.BinInfo.BuildGoVersion=${BuildGoVersion}' \
"

go build -ldflags "${LDFlags}"