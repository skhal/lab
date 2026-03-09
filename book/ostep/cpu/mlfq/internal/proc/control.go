// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/io"
)

// Cycler is the cpu.Clock interface used by [Process].
type Cycler interface {
	// Cycles returns current cpu.Cycle.
	Cycle() cpu.Cycle
}

// Control is the process controller. It provides API to change process state.
type Control struct {
	*Process

	clk Cycler
	io  *io.Runner
}

// Arrive marks the process arrive to the system.
func (ctl *Control) Arrive() {
	if ctl.arrive.set {
		return
	}
	ctl.state = StateReady
	ctl.arrive.set = true
	ctl.arrive.cycle = ctl.clk.Cycle()
}

// Run executes the process for one CPU cycle.
func (ctl *Control) Run() *io.Runner {
	switch ctl.state {
	case StateBlocked:
		panic(fmt.Errorf("%s: can't run", ctl))
	case StateReady:
		ctl.state = StateRunning
	case StateZombie:
		panic(fmt.Errorf("%s: can't run", ctl))
	}
	ctl.run()
	switch {
	case ctl.done():
		ctl.zombie()
	case ctl.ioCycle():
		ctl.callIO()
	}
	return ctl.io
}

func (ctl *Control) run() {
	if !ctl.firstRun.set {
		ctl.firstRun.set = true
		ctl.firstRun.cycle = ctl.clk.Cycle()
	}
	ctl.cycles++
}

func (ctl *Control) done() bool {
	return ctl.cycles == ctl.spec.CPUCycles
}

func (ctl *Control) zombie() {
	ctl.state = StateZombie
	if !ctl.complete.set {
		ctl.complete.set = true
		ctl.complete.cycle = ctl.clk.Cycle()
	}
}

func (ctl *Control) ioCycle() bool {
	if ctl.spec.IOAfterCPUCycles == 0 {
		return false
	}
	return ctl.cycles%ctl.spec.IOAfterCPUCycles == 0
}

func (ctl *Control) callIO() {
	ctl.state = StateBlocked
	ctl.io = io.NewRunner(ctl.spec.IOCycles)
}

// Update unblocks the process if it was blocked by IO that has completed. It
// should be called before the new cycle starts.
func (ctl *Control) Update() {
	switch ctl.state {
	case StateBlocked:
		if !ctl.io.Done() {
			return
		}
		ctl.state = StateReady
		ctl.io = nil
	}
}
