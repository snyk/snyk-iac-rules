#!/bin/bash

set -e

usage() {
    echo "$0 --tag=v0.0.0"
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

if [ -z "$TAG" ]; then
    usage
    exit 1
fi

cd $(dirname $0)/..

rm -rf dist

export VERSION="${TAG}"

echo "Updating NPM package version to ${TAG}"

mkdir -p dist/

for GOOS in linux darwin; do
    GOOS=$GOOS GOARCH=amd64 go build -a -o dist/snyk-iac-rules-$GOOS-x64 .
    GOOS=$GOOS GOARCH=arm64 go build -a -o dist/snyk-iac-rules-$GOOS-arm64 .
done
GOOS=windows GOARCH=amd64 go build -a -o dist/snyk-iac-rules-win.exe .

cp packaging/npm/passthrough.js dist/snyk-iac-rules
cp README.md dist/README.md
envsubst < packaging/npm/package.json.in > dist/package.json
