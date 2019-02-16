#!/usr/bin/env bash

set -e

VERSION="$1"
[[ -z "$VERSION" ]] && echo "No version defined" && exit -1

bash ./update_docs.bash
git comit -m "updating documentation for release of v$VERSION"

git tag "v$VERSION"
goreleaser --rm-dist