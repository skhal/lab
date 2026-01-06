// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Gotest runs `go test` on package for Go files.

Synopsis:

	check-go-test file [file ...]
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func main() {
	flag.Parse()
	if err := test.Run(flag.Args()...); err != nil {
		if !errors.Is(err, test.ErrTest) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
