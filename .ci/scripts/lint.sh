#!/bin/sh

set -e

go run golang.org/x/lint/golint@latest -set_exit_status ./...
