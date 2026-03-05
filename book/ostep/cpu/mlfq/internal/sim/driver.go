// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"iter"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
)

// Cycle is a single cycle.
type Cycle struct {
	// ID is the cycle identification.
	ID cpu.Cycle
}

// Policy is a scheduler interface used by simulator.
type Policy interface {
	// Add introduces a process to the policy.
	Add(policy.Process)

	// Next schedules a process to run. It ruturns true if scheduling succeeded,
	// else false.
	Next() bool

	// Process gives access to the selected process to run.
	Process() policy.Process
}

// Run drives processes pp with MLFQ policy using policy specifications pol.
// It returns a sequence of CPU cycles.
func Run(clk *cpu.Clock, pol Policy, pp []*proc.Process) iter.Seq[Cycle] {
	d := &driver{
		cpu:       clk,
		pol:       pol,
		processes: pp,
	}
	return d.Drive()
}

type driver struct {
	cpu *cpu.Clock
	pol Policy

	processes []*proc.Process
	completed []*proc.Process

	pending int // index of the next pending process
}

// Drive runs simulation and returns a sequence of CPU cycles.
func (dr *driver) Drive() iter.Seq[Cycle] {
	return func(yield func(Cycle) bool) {
		for dr.schedule() && yield(dr.cycle()) {
			dr.cpu.Next()
		}
	}
}

func (dr *driver) schedule() bool {
	// TODO(github.com/skhal/lab/issues/174): schedule pending processes with
	// matching Spec.Arrive cycle.
	return false
}

func (dr *driver) cycle() (cl Cycle) {
	defer func() {
		cl.ID = dr.cpu.Cycle()
	}()
	p := dr.pol.Process()
	if p == nil {
		return
	}
	p.(*proc.Process).Run()
	// TODO(github.com/skhal/lab/issues/174): handle process if completed
	return
}
