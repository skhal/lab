// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Allocator simulates memory allocation user library such as malloc(3). It
// demonstrates basics of managing free memory lists.
//
// SYNOPSIS
//
//	allocator [-base num] [-bounds num] [-n num] [-c noop|forward]
//
// EXAMPLE
//
//	% allocator/ -base 1024 -size 2048 -n 10 -c forward
//	configuration:
//	  base: 1024 size: 2048 coalesce: forward
//	  [1] free blocks 2046:1026
//
//	trace:
//	  malloc(225)
//	    [1] allocations 1026
//	    [1] free blocks 1819:1253
//	  malloc(326)
//	    [2] allocations 1026 1253
//	    [1] free blocks 1491:1581
//	  free(1026)
//	    [1] allocations 1253
//	    [2] free blocks 225:1026 1491:1581
//	  malloc(1892) malloc(1892): insufficient memory
//	  free(1253)
//	    [0] allocations
//	    [2] free blocks 225:1026 1819:1253
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
