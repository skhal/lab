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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/x/compdb/internal/compdb"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s package ...\n", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(targets []string) error {
	commands, err := compdb.GenCommands(targets)
	if err != nil {
		return err
	}
	return compdb.Print(os.Stdout, commands)
}
