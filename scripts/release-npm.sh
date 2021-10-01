#!/bin/bash

cd $(dirname $0)/..

rm -rf dist

VERSION="$1"
export VERSION="${VERSION}"

echo "Updating NPM package version to ${VERSION}"

mkdir -p dist/

for GOOS in linux darwin; do
    GOOS=$GOOS GOARCH=amd64 go build -a -o dist/snyk-iac-rules-$GOOS-amd64 .
    GOOS=$GOOS GOARCH=arm64 go build -a -o dist/snyk-iac-rules-$GOOS-arm64 .
done
GOOS=windows GOARCH=amd64 go build -a -o dist/snyk-iac-rules.exe .

cp packaging/npm/passthrough.js dist/snyk-iac-rules
envsubst < packaging/npm/package.json.in > dist/package.json
