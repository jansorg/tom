#!/usr/bin/env bash

go build ./gotime/docsGenerator/markdown \
    && ./markdown "$PWD/docs/markdown" \
    && rm -f markdown \
    && echo "Updated markdown documentation..." \
    || echo "Unable to update markdown documentation ..."

go build ./gotime/docsGenerator/man \
    && ./man "$PWD/docs/man" \
    && rm -f man \
    && echo "Updated man page documentation..." \
    || echo "Unable to update man page documentation ..."
