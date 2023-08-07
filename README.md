# go-windows

[![Build Status](https://beats-ci.elastic.co/job/Library/job/go-windows-mbp/job/main/badge/icon)](https://beats-ci.elastic.co/job/Library/job/go-windows-mbp/job/main/)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[godocs]:   http://godoc.org/github.com/elastic/go-windows

go-windows is a library for Go (golang) that provides wrappers to various
Windows APIs that are not covered by the stdlib or by
[golang.org/x/sys/windows](https://godoc.org/golang.org/x/sys/windows).

Goals / Features

- Does not use cgo.
- Provide abstractions to make using the APIs easier.
