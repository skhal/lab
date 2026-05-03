// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Htmlview serves a static HTML file.
//
// SYNOPSIS
//
//	htmlview [-http [host]:port] file
//
// # DESCRIPTION
//
// hemlview serves a static file at specified file.
//
// # EXAMPLE
//
// Serve a coverage report on the remote node:
//
//	remote % go test -count=1 -coverprofile=/tmp/sheet.out ./x/sheet/...
//	remote % go tool cover -html=/tmp/sheet.out -o /tmp/sheet_cover.html
//	remote % htmlview /tmp/sheet_cover.html
//
// Use port-forwarding to open the report on the client, e.g. laptop:
//
//	client % ssh -NL 8080:localhost:8080 remote
//	client % open https://localhost:8080/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var addr = flag.String("http", ":8080", "http server address")

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "%s file", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	file := flag.Arg(0)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, file)
	})
	fmt.Printf("serve on %s\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
