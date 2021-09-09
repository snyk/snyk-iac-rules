#!/bin/bash

GIT_COMMIT="$1"

NEW_TAG=$(./scripts/compute-release-tag.sh $GIT_COMMIT)

echo "Pushing ${VERSION} GitHub Release"

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
    echo "There is already a tag on this commit"
fi