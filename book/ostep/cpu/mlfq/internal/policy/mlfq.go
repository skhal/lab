// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/queue"
)

// Cycler is a CPU clock that gives access to current cycle.
type Cycler interface {
	// Cycle returns current CPU cycle.
	Cycle() cpu.Cycle
}

// New creates a MLFQ policy.
func New(spec Spec, c Cycler) *mlfq {
	queues := make([]*queue.RoundRobin, spec.NumQueues)
	for i := range spec.NumQueues {
		queues[i] = new(queue.RoundRobin)
	}
	return &mlfq{
		clk:    c,
		spec:   spec,
		queues: queues,
	}
}

// Process is a Process interface, used by MLFQ policy.
type Process interface {
	// ID returns process identifier.
	ID() int

	// Spec returns process specification.
	Spec() proc.Spec
}

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
	queues []*queue.RoundRobin

	last *process // last run job
}

// Add introduces a process to the scheduler.
func (pol *mlfq) Add(p Process) {
	pol.addToQueue(0, p)
}

func (pol *mlfq) addToQueue(q int, p Process) {
	pol.queues[q].Append(&process{
		proc: p,
		qid:  q,
	})
}

// Next picks up next process to run. It returns true if such process exists,
// else false.
func (pol *mlfq) Next() Process {
	pol.update()
	pol.last = pol.next()
	if pol.last == nil {
		return nil
	}
	return pol.last.proc
}

func (pol *mlfq) update() {
	if pol.last == nil {
		return
	}
	pol.last.cycles++
	switch pol.last.cycles {
	case pol.last.proc.Spec().CPUCycles:
		pol.remove(pol.last)
	case pol.spec.Allotment:
		pol.deprioritize(pol.last)
	}
}

func (pol *mlfq) remove(p *process) {
	if x := pol.queues[p.qid].Pop(); x.(*process) != p {
		panic(fmt.Errorf("remove: got %v, want %v", x.(*process), p))
	}
}

func (pol *mlfq) deprioritize(p *process) {
	if p.qid+1 == len(pol.queues) {
		// already lowest proiority
		return
	}
	pol.remove(p)
	pol.addToQueue(p.qid+1, p.proc)
}

func (pol *mlfq) next() *process {
	for _, q := range pol.queues {
		if q.Len() == 0 {
			continue
		}
		v := q.Next()
		return v.(*process)
	}
	return nil
}

type process struct {
	proc   Process
	qid    int
	cycles int
}

// String implements [fmt.Stringer] interface.
func (p *process) String() string {
	return fmt.Sprintf("pid:%d qid:%d cycles:%d %v", p.proc.ID(), p.qid, p.cycles, p.proc.Spec())
}
