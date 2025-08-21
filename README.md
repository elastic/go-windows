# go-windows

[![ci](https://github.com/elastic/go-windows/actions/workflows/ci.yml/badge.svg)](https://github.com/elastic/go-windows/actions/workflows/ci.yml)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[godocs]: https://pkg.go.dev/github.com/elastic/go-windows?GOOS=windows

go-windows is a library for Go (golang) that provides wrappers to various
Windows APIs that are not covered by the stdlib or by
[golang.org/x/sys/windows](https://godoc.org/golang.org/x/sys/windows).

## Goals / Features

- Does not use cgo.
- Provide abstractions to make using the APIs easier.

## Adding new Syscalls

Adding syscalls to zsyscall_windows.go is done by adding a comment under Syscalls in kernel32.go, psapi.go, etc describing the syscall.

Example:

```
// Syscalls
//sys   _GetPerformanceInfo(pi *PerformanceInformation, cb uint32) (err error) = psapi.GetPerformanceInfo
```

And then calling `go generate` as described in [doc.go](./doc.go)