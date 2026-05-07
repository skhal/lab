// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-go-test runs `go test` on package for Go files.
//
// SYNOPSIS
//
//	check-go-test file [file ...]
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
)

func main() {
	if err := check.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
