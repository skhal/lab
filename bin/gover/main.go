// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gover prints Go version used to build a Go binary. It pulls Go version from
// binary debug block with build information.
//
// SYNOPSIS
//
//	gover binary
package main

import (
	"debug/buildinfo"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(files []string) error {
	for _, file := range files {
		info, err := buildinfo.ReadFile(file)
		if err != nil {
			return err
		}
		fmt.Println(info.GoVersion, file)
	}
	return nil
}
