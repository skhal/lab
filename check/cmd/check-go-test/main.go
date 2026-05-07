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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s file ...\n", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if err := check.Run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
