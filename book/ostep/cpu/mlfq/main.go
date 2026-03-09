// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mlfq demonstrates Miltilevel Feedback Queue scheduler policy.
//
// SYNOPSIS
//
//	mlfq [-abort num] [-policy pol] [-proc spec[,:spec]]
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cmd"
)

func main() {
	if err := cmd.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
