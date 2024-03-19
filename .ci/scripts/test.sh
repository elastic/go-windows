#!/bin/sh

set -e

go mod verify
go run github.com/elastic/go-licenser@latest -d
out=$(go run golang.org/x/tools/cmd/goimports@latest -l -local github.com/elastic/go-windows .)

if [ ! -z "$out" ]; then
	printf "Run goimports on the code.\n"
	exit 1
fi

# Run the tests
export OUTPUT_JSON_FILE="build/test-report.out"
export OUTPUT_JUNIT_FILE="build/junit-${GO_VERSION}.xml"
mkdir -p build

go run gotest.tools/gotestsum@latest --no-color -f standard-quiet --jsonfile "$OUTPUT_JSON_FILE" --junitfile "$OUTPUT_JUNIT_FILE" ./...
