// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Mdview runs a server to preview a markdown file.

Synopsis:

	mdview file

The server listens on port :8080. It reads the markdown file from the file
system, render it to HTML, and serves the content.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/x/mdview/internal/mdview"
)

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "usage: %s file\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := mdview.Run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
