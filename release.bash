#!/usr/bin/env bash

set -e

VERSION="$1"
[[ -z "$VERSION" ]] && echo "No version defined" && exit 1

(go build . && ./tom completion > scripts/completions/tom.sh && chmod u+x scripts/completions/tom.sh)
git commit -m "prep: updating completions script for release of v$VERSION" . && git push

bash ./update_docs.bash
git commit -m "prep: updating documentation for release of v$VERSION" . && git push

# gotext generates data in unstable order, we can't use go generate here because it's called by goreleaser
#gotext -srclang=en update -out=catalog.go -lang=en,de
#git commit -m "prep: updating translations for release of v$VERSION" . && git push

git tag "v$VERSION"
#goreleaser --rm-dist
git push --tags
echo "The release will be built by GitHub Actions"