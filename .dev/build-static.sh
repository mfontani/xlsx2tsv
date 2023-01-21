#!/bin/sh

# fatal: detected dubious ownership in repository at ...
git config --global --add safe.directory "$(pwd)"
go build -tags timetzdata --ldflags "-X 'main.Version=$(git describe --tags)' -extldflags \"-static\" -s -w" -o xlsx2tsv .
