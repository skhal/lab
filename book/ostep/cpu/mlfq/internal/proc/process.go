// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

var lastID = 0

// Process is a process in the system.
type Process struct {
	id   int
	spec *Spec
}

// New creates a process with unique ID.
func New(s *Spec) *Process {
	lastID++
	return &Process{
		id:   lastID,
		spec: s,
	}
}

// ID returns process's identifier.
func (p *Process) ID() int {
	return p.id
}

// Run executes the process for one CPU cycle.
func (p *Process) Run() {
	// TODO(github.com/skhal/lab/issues/174): implement
}

// Spec gives access to the process's specification.
func (p *Process) Spec() Spec {
	return *p.spec
}
