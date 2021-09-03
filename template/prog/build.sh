#set -x

GitTag=$(git tag --sort=version:refname | tail -n 1)
GitCommitLog=$(git log --pretty=oneline -n 1)
GitCommitLog=${GitCommitLog//\'/\"}
BuildTime=$(date +'%Y.%m.%d.%H%M%S')
BuildGoVersion=$(go version)

LDFlags=" \
    -X 'main.GitTag=${GitTag}' \
    -X 'main.GitCommitLog=${GitCommitLog}' \
    -X 'main.BuildTime=${BuildTime}' \
    -X 'main.BuildGoVersion=${BuildGoVersion}' \
"

go build -ldflags "${LDFlags}"