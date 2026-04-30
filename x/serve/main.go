// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Serve serves files from a path via HTTP.
//
// SYNOPSIS
//
//	serve [-http addr] path
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var (
	addr = flag.String("http", ":8080", "local address for HTTP server")
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [flags] path\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return fmt.Errorf("missing path")
	}
	path := flag.Arg(0)
	if err := os.Chdir(path); err != nil {
		return err
	}
	fmt.Printf("serve %s at %s\n", path, *addr)
	return http.ListenAndServe(*addr, http.FileServer(http.Dir(".")))
}
