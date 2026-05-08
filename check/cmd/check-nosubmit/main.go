// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-nosubmit checks for presence of "DO NOT SUBMIT" comment.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s file ...\n", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if err := nosubmit.Run(flag.Args()...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
