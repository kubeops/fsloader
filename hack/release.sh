#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/fsloader"

export APPSCODE_ENV=prod

pushd $REPO_ROOT

rm -rf $REPO_ROOT/dist

./hack/make.py build
./hack/make.py push
./hack/make.py push

rm dist/.tag

popd
