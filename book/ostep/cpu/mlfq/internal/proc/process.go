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

// Cycler is the cpu.Clock interface used by [Process].
type Cycler interface {
	// Cycles returns current cpu.Cycle.
	Cycle() cpu.Cycle
}

// Process is a process in the system.
type Process struct {
	id   int
	spec Spec

	clk   Cycler
	state *state
}

type state struct {
	arrive struct {
		once  sync.Once
		cycle cpu.Cycle
	}

	firstRun struct {
		once  sync.Once
		cycle cpu.Cycle
	}

	complete struct {
		once  sync.Once
		cycle cpu.Cycle
	}

	cycles cpu.Cycle
}

// New creates a process with unique ID.
func New(s *Spec, clk Cycler) *Process {
	lastID++
	return &Process{
		id:    lastID,
		spec:  *s,
		clk:   clk,
		state: new(state),
	}
}

// ID returns process's identifier.
func (p *Process) ID() int {
	return p.id
}

// Arrive marks the process arrive to the system.
func (p *Process) Arrive() {
	p.state.arrive.once.Do(func() {
		p.state.arrive.cycle = p.clk.Cycle()
	})
}

// Run executes the process for one CPU cycle.
func (p *Process) Run() {
	p.state.firstRun.once.Do(func() {
		p.state.firstRun.cycle = p.clk.Cycle()
	})
	p.state.cycles++
	if p.Done() {
		p.state.complete.once.Do(func() {
			p.state.complete.cycle = p.clk.Cycle()
		})
	}
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

// Stat holds process metrics in cpu.Cycle units.
type Stat struct {
	// Response is the time between the process arrives and runs for the first
	// time.
	Response cpu.Cycle

	// Turnaround is the time between the process arrives and completes.
	Turnaround cpu.Cycle

	// Wait is the time the process spends not running on CPU.
	Wait cpu.Cycle
}

// Stat calculates process metrics.
func (p *Process) Stat() Stat {
	st := Stat{
		Response:   p.state.firstRun.cycle - p.state.arrive.cycle,
		Turnaround: p.state.complete.cycle - p.state.arrive.cycle,
	}
	st.Wait = st.Turnaround - p.state.cycles
	return st
}

// String implements [fmt.Stringer] interface.
func (p *Process) String() string {
	return fmt.Sprintf("pid:%d", p.id)
}
