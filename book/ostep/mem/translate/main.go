// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Translate demonstrates address translation using base and bounds CPU
// registers, where base is the offset of the virtual address space and bounds
// is the size of the virtual address space.
//
// SYNOPSIS
//
//	translate [-base num] [-bounds num]
//
// # DESCRIPTION
//
// Address translation uses base to get physical address:
//
//	phys = virt + base
//
// and verifies that the address is within boundaries using bounds.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := new(command).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
