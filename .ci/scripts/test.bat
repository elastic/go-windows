
go mod verify || EXIT /B 1
go run github.com/elastic/go-licenser@latest -d || EXIT /B 1

go run golang.org/x/tools/cmd/goimports@latest -l -local github.com/elastic/go-windows . || EXIT /B 1

SET OUTPUT_JSON_FILE=build\output-report.out
SET OUTPUT_JUNIT_FILE=build\junit-%GO_VERSION%.xml

go run gotest.tools/gotestsum@latest --no-color -f standard-quiet --jsonfile "$OUTPUT_JSON_FILE" --junitfile "$OUTPUT_JUNIT_FILE" ./...
