// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Gotest runs `go test` on package for Go files.

Synopsis:

  gotest file [file ...]
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/gotest/internal/test"
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
