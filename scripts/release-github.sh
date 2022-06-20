#!/usr/bin/env bash

set -eo pipefail

if ! which goreleaser >/dev/null ; then
    go install github.com/goreleaser/goreleaser@v1.9.2
fi

if ! which go-licenses >/dev/null ; then
    go install github.com/google/go-licenses@latest
fi

# Generate acknowledgements
go-licenses save . --save_path=./acknowledgements
tar -cvf ./acknowledgements.tar.gz -C ./acknowledgements .
rm -rf ./acknowledgements

# Check configuration
goreleaser check

FLAGS=""
FLAGS+="--rm-dist "

# Only CI system should publish artifacts
if [ "$CI" != true ]; then
    FLAGS+="--skip-announce "
    FLAGS+="--skip-publish "
    FLAGS+="--snapshot "
fi

CMD="goreleaser release ${FLAGS}"

echo "+ Using goreleaser"
echo "+ CMD=${CMD}"

$CMD
