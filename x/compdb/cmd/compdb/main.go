// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Compdb generates a compilation database from a Bazel project.
//
// Synopsis
//
//	compdb TARGET ...
//
// LLVM defines a compilation database as a set of instructions to build
// source code in the project. Each instruction is a build command for a
// single target, https://clang.llvm.org/docs/JSONCompilationDatabase.html.
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/x/compdb/internal/compdb"
)

func main() {
	if err := compdb.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
