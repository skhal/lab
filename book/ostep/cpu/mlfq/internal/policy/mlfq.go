// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import "github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"

// Cycler is a CPU clock that gives access to current cycle.
type Cycler interface {
	// Cycle returns current CPU cycle.
	Cycle() cpu.Cycle
}

// New creates a MLFQ policy.
func New(spec Spec, c Cycler) *mlfq {
	return &mlfq{
		clk:    c,
		spec:   spec,
		queues: make([][]Process, spec.NumQueues),
	}
}

// Process is a Process interface, used by MLFQ policy.
type Process any

// mlfq implements Multilevel Feedback Queue scheduling policy. It uses the
// following rules:
//
//  1. Add new jobs to the top-priority queue
//  2. Round-robin jobs from the highest priority non-empty queue
//  3. Decrease job priority if it used allotted CPU time share
//  4. Reset priorities once in a while
type mlfq struct {
	clk    Cycler
	spec   Spec
	queues [][]Process
	proc   Process // last run job
}

// Add introduces a process to the scheduler.
func (pol *mlfq) Add(j Process) {
	pol.queues[0] = append(pol.queues[0], j)
}

// Next picks up next process to run. It returns true if such process exists,
// else false.
func (pol *mlfq) Next() bool {
	pol.update()
	pol.next()
	return pol.proc != nil
}

// Process gives access to currently selected process.
func (pol *mlfq) Process() Process {
	return pol.proc
}

func (pol *mlfq) update() {
	// TODO(github.com/skhal/lab/issues/174): update the last job based on the
	// last run operation
}

func (pol *mlfq) next() {
	// TODO(github.com/skhal/lab/issues/174): pick up next job
	// TODO(github.com/skhal/lab/issues/174): record the next job
}
