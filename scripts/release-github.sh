#!/usr/bin/env bash

set -eo pipefail

if ! which goreleaser >/dev/null ; then
    go install github.com/goreleaser/goreleaser/v2@latest
fi

# Check configuration
goreleaser check

FLAGS=""
FLAGS+="--clean "

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
