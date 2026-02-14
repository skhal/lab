// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-pbdoc verifies that every item in a Protobuf definition includes a
// documentation comment.
//
// Synopsis:
//
//	check-proto path/to/file.proto
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-pbdoc/internal/check"
)

func main() {
	if err := check.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
