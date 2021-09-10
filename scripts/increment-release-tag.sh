#!/bin/bash

GIT_COMMIT="$1"

LATEST_TAG_WITH_V=$(git describe --tags $(git rev-list --tags --max-count=1))
LATEST_TAG=${LATEST_TAG_WITH_V:1}

IFS=. SPLIT_LATEST_TAG=(${LATEST_TAG##*-})
MAJOR=${SPLIT_LATEST_TAG[0]}
MINOR=${SPLIT_LATEST_TAG[1]}
PATCH=${SPLIT_LATEST_TAG[2]}

PATCH=$((PATCH+1))

NEW_TAG="v$MAJOR.$MINOR.$PATCH"

echo "Updating $LATEST_TAG to $NEW_TAG"

# Get current hash and see if it already has a tag
NEEDS_TAG=$(git describe --contains $GIT_COMMIT 2>/dev/null)
if [ -z "$NEEDS_TAG" ]; then
    GIT_EMAIL=$(git show -s --format='%ae' $GIT_COMMIT) && \
    echo "Email used for the commit was ${GIT_EMAIL}" && \
    git config credential.helper 'cache --timeout=120' && \
    git config user.email "${GIT_EMAIL}" && \
    git config user.name "Automated Release" && \
    echo "git tag -a ${NEW_TAG} -m \"Release ${NEW_TAG}\"" && \
    git tag -a "${NEW_TAG}" -m "Release ${NEW_TAG}" && \
    echo "Tagged commit with $NEW_TAG" && \
    echo "git push origin ${NEW_TAG}"
    git push origin "${NEW_TAG}" && \
    echo "New tag pushed to GitHub"
else
    echo "There already is a tag on this commit"
fi