// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"fmt"
	"iter"
	"slices"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/io"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
)

var abortCycle cpu.Cycle = 1000

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
func Run(clk *cpu.Clock, pol Policy, specs []*proc.Spec) ([]*proc.Process, iter.Seq[Cycle]) {
	var (
		pp   = make([]*proc.Process, 0, len(specs))
		ctls = make([]*proc.Control, 0, len(specs))
	)
	for _, spec := range specs {
		p, c := proc.New(spec, clk)
		pp = append(pp, p)
		ctls = append(ctls, c)
	}
	d := &driver{
		cpu:       clk,
		pol:       pol,
		processes: ctls,
	}
	return pp, d.Drive()
}

type driver struct {
	cpu *cpu.Clock
	pol Policy

	// processes are all processes in the system: pending, running, or completed.
	processes []*proc.Control
	pending   int // index of the next pending process
	completed int

	cycle Cycle

	// io stores processes, that are blocked on the IO, unblocked processes,
	// and IO request by the last process. Notes:
	//
	// 1. the driver calls processes, whose IO completed, to update the state,
	//    i.e., unblock. It happens before the cycle starts using [io.unblocked]
	//    slice.
	//
	// 2. separate cache from blocked to avoid running CPU and IO for a single
	//    process in the same cycle: runCPU makes a request for IO, runIO
	//    emulates the IO. These should happen in cycles N and N+1.
	//
	// 3. runIO executed IO cycles and cleans processes whose IO finished. We
	//    want to keep these for the beginning of the next cycle to update the
	//    processes.
	io struct {
		blocked   []*ioblock // processes blocked on IO
		unblocked []*ioblock // processes with completed IO, pre-update
		cache     *ioblock   // last process IO request, pre-run IO
	}
}

type ioblock struct {
	proc *proc.Control
	io   *io.Runner
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
	if clk := dr.cpu.Cycle(); clk == abortCycle {
		panic(fmt.Errorf("clk %d: abort", clk))
	}
	if dr.completed == len(dr.processes) {
		return false
	}
	dr.unblock()
	dr.schedule()
	dr.cpu.Next()
	dr.run()
	return true
}

func (dr *driver) unblock() {
	for _, b := range dr.io.unblocked {
		b.proc.Update()
	}
	dr.io.unblocked = nil
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
	dr.runCPU()
	dr.runIO()
}

func (dr *driver) runCPU() {
	x, pri := dr.pol.Next()
	if x == nil {
		return
	}
	p := x.(*proc.Control)
	switch io := p.Run(); {
	case io != nil:
		dr.io.cache = &ioblock{
			proc: p,
			io:   io,
		}
	case p.Done():
		dr.completed++
	}
	dr.cycle.Proc = p.Process
	dr.cycle.Priority = pri
}

func (dr *driver) runIO() {
	for _, b := range dr.io.blocked {
		b.io.Run()
	}
	dr.io.blocked = slices.DeleteFunc(dr.io.blocked, func(b *ioblock) bool {
		if !b.io.Done() {
			return false
		}
		dr.io.unblocked = append(dr.io.unblocked, b)
		return true
	})
	if dr.io.cache != nil {
		dr.io.blocked = append(dr.io.blocked, dr.io.cache)
		dr.io.cache = nil
	}
}
