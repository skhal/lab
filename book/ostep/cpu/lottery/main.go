// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Lottery runs simulation of a lottery scheduler. A lottery scheduler uses
// job weights to randomly pick up the next job to run.
//
// This implementation uses a concept of tickets for weights:
//
//   - each job gets Ni tickets
//   - there are total N tickets that is equal to the sum of all tickets across
//     jobs: sum(Ni)
//
// Random weight is then Ni/sum(Ni).
//
// SYNOPSIS
//
//	lottery [-jobs spec[,:spec]]
//
// where `spec` is a colon-separated list of [job.Spec] fields: length,
// tickets.
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
