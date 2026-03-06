// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"
	"sync"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

var lastID = 0

// Process is a process in the system.
type Process struct {
	id   int
	spec Spec

	state *state
}

type state struct {
	once   sync.Once
	arrive cpu.Cycle

	cycles cpu.Cycle
}

// New creates a process with unique ID.
func New(s *Spec) *Process {
	lastID++
	return &Process{
		id:    lastID,
		spec:  *s,
		state: new(state),
	}
}

// ID returns process's identifier.
func (p *Process) ID() int {
	return p.id
}

// Arrive marks the process arrive to the system.
func (p *Process) Arrive(c cpu.Cycle) {
	p.state.once.Do(func() {
		p.state.arrive = c
	})
}

// Run executes the process for one CPU cycle.
func (p *Process) Run() {
	p.state.cycles++
}

// Done reports whether the process completed.
func (p *Process) Done() bool {
	return p.state.cycles == cpu.Cycle(p.spec.CPUCycles)
}

// Cycles returns the number of completed CPU cycles.
func (p *Process) Cycles() cpu.Cycle {
	return p.state.cycles
}

// Spec gives access to the process's specification.
func (p *Process) Spec() Spec {
	return p.spec
}

// String implements [fmt.Stringer] interface.
func (p *Process) String() string {
	return fmt.Sprintf("pid:%d", p.id)
}
