// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Mdview runs a server to preview markdown files rotted at current working
directory.

Synopsis:

	mdview

The server listens on port :8080. Make a GET HTTP request to render a markdown
file given by URL path:

	open localhost:8080/path/to/README.md

The server reads the file on every request.
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
		fmt.Fprintf(out, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	if err := mdview.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
