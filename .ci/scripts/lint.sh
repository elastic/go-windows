#!/bin/sh

set -e

go run honnef.co/go/tools/cmd/staticcheck@2022.1 ./...
