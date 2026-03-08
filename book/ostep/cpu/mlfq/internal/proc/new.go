// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

var lastID = 0

// New creates a process with unique ID.
func New(s *Spec, clk Cycler) (*Process, *Control) {
	lastID++
	p := &Process{
		id:    lastID,
		spec:  *s,
		state: new(state),
	}
	c := &Control{
		Process: p,
		clk:     clk,
	}
	return p, c
}
