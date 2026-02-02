// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-nextid verifies that a Protocol buffer definition has a next-id
// comment set to the next value to the maximum used field identifier.
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-nextid/internal/nextid"
)

func main() {
	if err := nextid.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
