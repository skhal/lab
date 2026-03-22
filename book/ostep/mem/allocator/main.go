// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Allocator simulates memory allocation user library such as malloc(3). It
// demonstrates basics of managing free memory lists.
//
// SYNOPSIS
//
//	allocator [-base num] [-bounds num]
//
// EXAMPLE
//
//	% allocator -base 1024 -size 2048
//	configuration:
//	  base: 1024 size: 2048 free[1] 2044:1028
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
