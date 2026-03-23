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
//	 % allocator/ -base 1024 -size 2048 -n 10        configuration:
//		 base: 1024 size: 2048
//		 [1] free blocks 2046:1026
//
//	 trace:
//		 malloc(1503)
//			 [1] allocations 1026
//			 [1] free blocks 541:2531
//		 free(1026)
//			 [0] allocations
//			 [2] free blocks 1503:1026 541:2531
//		 malloc(1937) malloc(1937): insufficient memory
//		 malloc(1646) malloc(1646): insufficient memory
//		 malloc(498)
//			 [1] allocations 1026
//			 [2] free blocks 1003:1526 541:2531
//		 free(1026)
//			 [0] allocations
//			 [3] free blocks 498:1026 1003:1526 541:2531
//		 malloc(1297) malloc(1297): insufficient memory
//		 malloc(1379) malloc(1379): insufficient memory
//		 malloc(772)
//			 [1] allocations 1526
//			 [3] free blocks 498:1026 229:2300 541:2531
//		 malloc(1754) malloc(1754): insufficient memory
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
