// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Kscope parses a file using "kaleidoscope" syntax and prints parsed source.
//
// SYNOPSIS
//
//	kscope file.ks
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/x/kscope/internal/ast"
	"github.com/skhal/lab/x/kscope/internal/parse"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s file\n", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fname string) error {
	n, err := parseFile(fname)
	if err != nil {
		return err
	}
	fmt.Println(n)
	return nil
}

func parseFile(name string) (ast.Node, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	n, err := parse.Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}
