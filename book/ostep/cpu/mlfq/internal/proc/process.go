// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

// Process is a process in the system.
type Process struct {
	id   int
	spec Spec

	arrive struct {
		set   bool
		cycle cpu.Cycle
	}

	firstRun struct {
		set   bool
		cycle cpu.Cycle
	}

	complete struct {
		set   bool
		cycle cpu.Cycle
	}

	cycles cpu.Cycle

	state State
}

// Blocked checks whether the process state is [StateBlocked].
func (p *Process) Blocked() bool {
	return p.state == StateBlocked
}

// Cycles returns the number of completed CPU cycles.
func (p *Process) Cycles() cpu.Cycle {
	return p.cycles
}

// Done reports whether the process completed.
func (p *Process) Done() bool {
	return p.state == StateZombie
}

// ID returns process's identifier.
func (p *Process) ID() int {
	return p.id
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
		Response:   p.firstRun.cycle - p.arrive.cycle,
		Turnaround: p.complete.cycle - p.arrive.cycle,
	}
	st.Wait = st.Turnaround - p.cycles
	return st
}

// State returns current process state.
func (p *Process) State() State {
	return p.state
}

// String implements [fmt.Stringer] interface.
func (p *Process) String() string {
	return fmt.Sprintf("pid:%d [%s]", p.id, p.state)
}
