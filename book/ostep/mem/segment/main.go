// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Segment demonstrates address translation with support for memory segments.
//
// # SYNOPSIS
//
//	segment [-n num] [-segA base:bounds] [-segB base:bounds] [-vm-bounds num]
//
// # DESCRIPTION
//
// A memory segment is a logical block of virtual memory. It has own base and
// bounds settings, along with permission bits (rwx), direction (positive or
// negative offset).
//
// Segment tool uses multiple segments to run address translation of virtual
// to physical address, referenced within different segments.
//
// EXAMPLE
//
//	 % segment -segA 2:2 -segB 4:2
//	 configuration:
//		 virtual address bounds: 4KB
//		 SEG0 base: 2KB bounds: 2KB dir: positive virt-base: 0KB
//		 SEG1 base: 4KB bounds: 2KB dir: negative virt-base: 4KB
//
//	 translations:
//		 virt: 3462 (SEG1) phys: 5510
//		 virt: 3046 (SEG0) segmentation fault
//		 virt: 2484 (SEG0) segmentation fault
//		 virt: 2949 (SEG0) segmentation fault
//		 virt: 342 (SEG0) phys: 2390
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
