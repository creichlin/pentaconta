#!/usr/bin/env bash

cd "$(dirname $(dirname "$0"))"

export GOPATH=$(dirname $(dirname $(dirname $(dirname $(pwd)))))
go install -v gitlab.com/creichlin/pentaconta/...
export PATH=$PATH:$GOPATH/bin
go test gitlab.com/creichlin/pentaconta/...
