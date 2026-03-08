// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import "github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"

// Cycler is the cpu.Clock interface used by [Process].
type Cycler interface {
	// Cycles returns current cpu.Cycle.
	Cycle() cpu.Cycle
}

// Control is the process controller. It provides API to change process state.
type Control struct {
	*Process
	clk Cycler
}

// Arrive marks the process arrive to the system.
func (ctl *Control) Arrive() {
	if ctl.arrive.set {
		return
	}
	ctl.arrive.set = true
	ctl.arrive.cycle = ctl.clk.Cycle()
}

// Run executes the process for one CPU cycle.
func (ctl *Control) Run() {
	if !ctl.firstRun.set {
		ctl.firstRun.set = true
		ctl.firstRun.cycle = ctl.clk.Cycle()
	}
	ctl.cycles++
	if ctl.Done() {
		if !ctl.complete.set {
			ctl.complete.set = true
			ctl.complete.cycle = ctl.clk.Cycle()
		}
	}
}
