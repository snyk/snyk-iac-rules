#!/bin/bash

cd $(dirname $0)/..

rm -rf dist

GIT_COMMIT="$1"

VERSION=$(./scripts/compute-release-tag.sh $GIT_COMMIT)
export VERSION="${VERSION}"

echo "Updating NPM package version to ${VERSION}"

mkdir -p dist/

for GOOS in linux darwin; do
    GOOS=$GOOS GOARCH=amd64 go build -a -o dist/snyk-iac-custom-rules-$GOOS-amd64 .
done
GOOS=windows GOARCH=amd64 go build -a -o dist/snyk-iac-custom-rules.exe .

cp packaging/npm/passthrough.js dist/snyk-iac-custom-rules
envsubst < packaging/npm/package.json.in > dist/package.json