// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmd provides basic constructs to run the lottery scheduler, i.e.,
// a runnable command. The command parses the flags, runs the simulation, and
// prints a report.
//
// SYNOPSIS
//
//	if err := cmd.Run(); err != nil {
//		log.Fatal(err)
//	}
package cmd
