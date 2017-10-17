#!/usr/bin/env bash

pushd $GOPATH/src/github.com/appscode/fsloader/hack/gendocs
go run main.go
popd
