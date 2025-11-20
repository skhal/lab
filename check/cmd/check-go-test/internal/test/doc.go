// Copyright 2025 Samvel Khalatyan. All rights reserved.

// Package test gives access to failed `go test` for Go file packages.
//
// Given a list of files, it collects packages for Go files, and runs `go test`
// in JSON output mode. It keeps track of test outputs and provides a quick
// access to failed test outputs.
package test
