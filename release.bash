#!/usr/bin/env bash

set -e

VERSION="$1"
[[ -z "$VERSION" ]] && echo "No version defined" && exit -1

bash ./update_docs.bash
git tag "v$VERSION"
goreleaser --rm-dist