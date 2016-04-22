#!/usr/bin/env bash
. ../etc/env.conf

go generate internal/version
go build -o $GOPATH/bin/$PROJNAME -race main
go list ./... | grep -v vendor/ | xargs -L1 go fmt
