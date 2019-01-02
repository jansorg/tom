#!/usr/bin/env bash

go build ./gotime/docsGenerator/markdown \
    && ./markdown \
    && rm -f markdown \
    && echo "Updated markdown documentation..." \
    || echo "Unable to update markdown documentation ..."
