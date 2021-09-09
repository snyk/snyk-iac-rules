#!/bin/bash

LATEST_TAG_WITH_V=$(git describe --tags $(git rev-list --tags --max-count=1))
LATEST_TAG=${LATEST_TAG_WITH_V:1}

IFS=. SPLIT_LATEST_TAG=(${LATEST_TAG##*-})
MAJOR=${SPLIT_LATEST_TAG[0]}
MINOR=${SPLIT_LATEST_TAG[1]}
PATCH=${SPLIT_LATEST_TAG[2]}

PATCH=$((PATCH+1))

NEW_TAG="v$MAJOR.$MINOR.$PATCH"

echo "${NEW_TAG}"