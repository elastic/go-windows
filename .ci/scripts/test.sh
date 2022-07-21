#!/bin/sh

set -e

go mod verify
go run github.com/elastic/go-licenser@latest -d
sh .ci/scripts/format.sh
sh .ci/scripts/lint.sh

# Run the tests
export OUTPUT_JSON_FILE="build/test-report.out"
export OUTPUT_JUNIT_FILE="build/junit-${GO_VERSION}.xml"
mkdir -p build

go run gotest.tools/gotestsum@latest --no-color -f standard-quiet --jsonfile "$OUTPUT_JSON_FILE" --junitfile "$OUTPUT_JUNIT_FILE" ./...
