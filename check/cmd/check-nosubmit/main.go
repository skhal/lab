// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-nosubmit checks for presence of "DO NOT SUBMIT" comment.
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

func main() {
	cfg := &nosubmit.Config{
		ReadFileFn: os.ReadFile,
	}
	if err := nosubmit.Run(cfg, os.Args[1:]...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
