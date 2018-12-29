#!/usr/bin/env bash

go build ./gotime/docs/markdown
./markdown
rm -f markdown