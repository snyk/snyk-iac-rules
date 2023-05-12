#!/usr/bin/env bash

set -eo pipefail

if ! which svu >/dev/null ; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
    brew install caarlos0/tap/svu@1.7.0
  else
    echo 'deb [trusted=yes] https://apt.fury.io/caarlos0/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
    sudo apt update
    sudo apt install svu=1.7.0
  fi
fi

TAG=$(svu)

if [ $(git tag -l "TAG") ]; then
    echo "Tag already exists!"
    exit 0
else
    # Hard-code user config
    git config user.email "cloud-dev-ex@users.noreply.github.com"
    git config user.name "Snyk"

    # Push new tag
    git tag -a "${TAG}" -m "Release ${TAG}"
    git push origin "${TAG}"

    # Set environment variable
    echo "export TAG=${TAG}" >> $BASH_ENV
fi
