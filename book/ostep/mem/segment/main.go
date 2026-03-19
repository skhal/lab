// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Segment demonstrates address translation with support for memory segments.
//
// # SYNOPSIS
//
//	segment
//
// # DESCRIPTION
//
// A memory segment is a logical block of virtual memory. It has own base and
// bounds settings, along with permission bits (rwx), direction (positive or
// negative offset).
//
// Segment tool uses multiple segments to run address translation of virtual
// to physical address, referenced within different segments.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := newCommand().Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
