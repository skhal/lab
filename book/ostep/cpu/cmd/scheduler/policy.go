// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Policy enumerates available scheduling policies.
//
//go:generate stringer -type=Policy -linecomment
type Policy int

const (
	_ Policy = iota
	// PolicyFIFO runs first-in-first-out job.
	PolicyFIFO // fifo
	// PolicySJF runs the job that is shortest to finish.
	PolicySJF // sjf
	// PolicySTCF preempts currently running job to pick up the shortest to
	// complete job.
	PolicySTCF // stcf
)
