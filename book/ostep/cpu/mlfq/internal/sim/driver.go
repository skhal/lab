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

// Cycle is a single CPU cycle.
type Cycle struct {
	// ID is the cycle identification.
	ID cpu.Cycle

	// Proc is a process run in a given cycle. It might be nil if no process was
	// scheduled, e.g., a process didn't arrive in the system yet and CPU sits
	// idle.
	Proc *proc.Process

	// Priority is the process priority.
	Priority policy.Priority
}

// Policy is a scheduler interface used by simulator.
type Policy interface {
	// Add introduces a process to the policy.
	Add(policy.Process)

	// Next schedules a process to run. It ruturns true if scheduling succeeded,
	// else false.
	Next() (policy.Process, policy.Priority)
}

// Run drives processes pp with MLFQ policy using policy specifications pol.
// It returns a sequence of CPU cycles.
func Run(clk *cpu.Clock, pol Policy, procs []*proc.Process) iter.Seq[Cycle] {
	d := &driver{
		cpu:       clk,
		pol:       pol,
		processes: procs,
	}
	return d.Drive()
}

type driver struct {
	cpu *cpu.Clock
	pol Policy

	// processes are all processes in the system: pending, running, or completed.
	processes []*proc.Process
	pending   int             // index of the next pending process
	completed []*proc.Process // completed processes

	cycle Cycle
}

// Drive runs simulation and returns a sequence of CPU cycles.
func (dr *driver) Drive() iter.Seq[Cycle] {
	return func(yield func(Cycle) bool) {
		for dr.next() && yield(dr.cycle) {
			continue
		}
	}
}

func (dr *driver) next() bool {
	if len(dr.completed) == len(dr.processes) {
		return false
	}
	dr.schedule()
	dr.cpu.Next()
	dr.run()
	return true
}

func (dr *driver) schedule() {
	for _, proc := range dr.processes[dr.pending:] {
		if proc.Spec().Arrive != dr.cpu.Cycle() {
			break
		}
		dr.pol.Add(proc)
		dr.pending++
	}
}

func (dr *driver) run() {
	dr.cycle = Cycle{
		ID: dr.cpu.Cycle(),
	}
	x, pri := dr.pol.Next()
	if x == nil {
		return
	}
	p := x.(*proc.Process)
	p.Run()
	if p.Done() {
		dr.completed = append(dr.completed, p)
	}
	dr.cycle.Proc = p
	dr.cycle.Priority = pri
}
