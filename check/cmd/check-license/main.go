// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Check-copyright verifies that the file includes a copyright statement.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-license/internal/license"
)

const argsLen = 1

func init() {
	flag.Usage = func() {
		header := func() string {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "Usage: %s file [file ...]\n", flag.CommandLine.Name())
			return buf.String()
		}
		fmt.Fprint(flag.CommandLine.Output(), header())
		flag.PrintDefaults()
	}
}

func main() {
	opts := new(license.RunOptions)
	opts.RegisterFlags(flag.CommandLine)
	flag.Parse()
	if err := license.Run(flag.Args(), opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
