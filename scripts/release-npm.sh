#!/bin/bash

set -e

usage() {
    echo "Incorrect usage: $0 --tag=v0.0.0"
}

for i in "$@"; do
    case $i in
    --tag=*)
        TAG="${i#*=}"
        shift
        ;;
    *)
        usage
        exit 1
        ;;
    esac
done

if ! which goreleaser >/dev/null ; then
    go install github.com/goreleaser/goreleaser/v2@latest
fi

# Check configuration
goreleaser check

# Override tag for GoReleaser so it uses the one provided in the flag
export GORELEASER_CURRENT_TAG="${TAG}"

CMD="goreleaser build --snapshot --clean"

echo "+ Using goreleaser"
echo "+ CMD=${CMD}"

$CMD

echo "Updating NPM package version to ${TAG}"

cp packaging/npm/passthrough.js dist/snyk-iac-rules
cp README.md dist/README.md
# Use the tag provided in the flag for the version field in the package.json
export VERSION="${TAG}"
envsubst < packaging/npm/package.json.in > dist/package.json
