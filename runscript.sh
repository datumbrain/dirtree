#!/bin/zsh

BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse HEAD)"
#VERSION="$(git describe --tags --abbrev=0 | tr -d "\n")"


go build -o bin/main -ldflags="-X 'github.com/dirtree/internal/version.buildDate=${BUILD_DATE}' -X 'github.com/dirtree/internal/version.gitCommit=${GIT_COMMIT}'" *.go
