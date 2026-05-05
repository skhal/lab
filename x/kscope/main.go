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

	"github.com/skhal/lab/x/kscope/internal/parse"
)

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
	f, err := parse.ParseFile(fname)
	if err != nil {
		return err
	}
	fmt.Println(f)
	return nil
}
