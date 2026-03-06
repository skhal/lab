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
		spec:   spec,
		clk:    c,
		queues: queues,
	}
}

// Process is a Process interface, used by MLFQ policy.
type Process interface {
	// ID returns process identifier.
	ID() int

	// Spec returns process specification.
	Spec() proc.Spec

	// Done returns true if the process completed, else false.
	Done() bool
}

// mlfq implements Multilevel Feedback Queue scheduling policy. It uses the
// following rules:
//
//  1. Add new processes to the top-priority queue
//  2. Round-robin processes from the highest priority non-empty queue
//  3. Decrease process priority if it used allotted CPU time share
//  4. Reset priorities once in a while
type mlfq struct {
	spec Spec

	clk    Cycler
	queues []*queue.RoundRobin
	last   *process // last run process
}

// Add injects the new process to the highest priority queue.
func (pol *mlfq) Add(p Process) {
	pol.addToQueue(0, p)
}

func (pol *mlfq) addToQueue(q int, p Process) {
	pol.queues[q].Append(&process{
		proc: p,
		qid:  q,
	})
}

// Priority is the process priority.
type Priority int

// Next picks up next process to run and returns it along with process's
// priority. It returns a nil process and undefined priority if the scheduler
// can't pick up the next process, e.g. there are no processes available to the
// scheduler.
func (pol *mlfq) Next() (Process, Priority) {
	pol.update()
	pol.last = pol.next()
	if pol.last == nil {
		return nil, 0
	}
	return pol.last.proc, Priority(pol.last.qid)
}

func (pol *mlfq) update() {
	if pol.last == nil {
		return
	}
	pol.last.cycles++
	if pol.last.proc.Done() {
		pol.remove(pol.last)
		return
	}
	if pol.last.cycles == pol.spec.Allotment {
		pol.deprioritize(pol.last)
	}
}

func (pol *mlfq) remove(p *process) {
	if x := pol.queues[p.qid].Pop(); x == nil {
		panic(fmt.Errorf("remove: failed to pop %v", p))
	} else if x.(*process) != p {
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
	cycles cpu.Cycle
}

// String implements [fmt.Stringer] interface.
func (p *process) String() string {
	return fmt.Sprintf("pid:%d qid:%d cycles:%d %v", p.proc.ID(), p.qid, p.cycles, p.proc.Spec())
}
