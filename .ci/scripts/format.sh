#!/bin/sh

out=$(go run golang.org/x/tools/cmd/goimports@latest -l -local github.com/elastic/go-windows .)

if [ ! -z "$out" ]; then
	printf "Run goimports on the code.\n" 
	exit 1
fi
