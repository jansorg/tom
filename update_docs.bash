#!/usr/bin/env bash

go build ./go-tom/docsGenerator/markdown \
    && git rm --cached docs/markdown/* \
    && ./markdown "$PWD/docs/markdown" \
    && rm -f markdown \
    && echo "Updated markdown documentation..." && ls && git add docs/markdown \
    || echo "Unable to update markdown documentation ..."

go build ./go-tom/docsGenerator/man \
    && git rm --cached docs/man/* \
    && ./man "$PWD/docs/man" \
    && rm -f man \
    && echo "Updated man page documentation..." && git add docs/man \
    || echo "Unable to update man page documentation ..."
