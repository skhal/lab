// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"fmt"
	"iter"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/queue"
)

// Priority is the process priority.
type Priority int

const maxPriority = Priority(0)

// Cycler is a CPU clock that gives access to current cycle.
type Cycler interface {
	// Cycle returns current CPU cycle.
	Cycle() cpu.Cycle
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

// New creates a MLFQ policy.
func New(spec Spec, c Cycler) *mlfq {
	queues := make([]*queue.RoundRobin, spec.Priorities)
	for i := range spec.Priorities {
		queues[i] = new(queue.RoundRobin)
	}
	return &mlfq{
		spec:   spec,
		clk:    c,
		queues: queues,
	}
}

// Add injects the new process to the highest priority queue.
func (pol *mlfq) Add(p Process) {
	p.Arrive()
	pol.addToQueue(maxPriority, p)
}

func (pol *mlfq) addToQueue(prio Priority, p Process) {
	pol.queues[prio].Append(&process{
		proc: p,
		prio: prio,
	})
}

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
	return pol.last.proc, Priority(pol.last.prio)
}

func (pol *mlfq) update() {
	if pol.last == nil {
		return
	}
	pol.last.cycles++
	switch {
	case pol.last.proc.Done():
		pol.remove(pol.last)
	case pol.last.atAllotment(pol.spec.Allotment):
		pol.deprioritize(pol.last)
	}
	if pol.canBoost() {
		pol.boost()
	}
}

func (pol *mlfq) canBoost() bool {
	return pol.clk.Cycle()%pol.spec.BoostCycles == 0
}

func (pol *mlfq) remove(p *process) {
	if x := pol.queues[p.prio].Pop(); x == nil {
		panic(fmt.Errorf("remove: failed to pop %v", p))
	} else if x.(*process) != p {
		panic(fmt.Errorf("remove: got %v, want %v", x.(*process), p))
	}
}

func (pol *mlfq) deprioritize(p *process) {
	if int(p.prio+1) == len(pol.queues) {
		// already lowest proiority
		return
	}
	pol.remove(p)
	pol.addToQueue(p.prio+1, p.proc)
}

func (pol *mlfq) boost() {
	// Move processes from lower priorities to the highest priority queue.
	// Therefore exclude highest priority queue.
	// Process lowest priority queues first to give these processes a change
	// to run first.
	for q := range reverse(pol.queues[1:]) {
		pol.prioritize(q)
	}
}

func reverse(qq []*queue.RoundRobin) iter.Seq[*queue.RoundRobin] {
	return func(yield func(*queue.RoundRobin) bool) {
		for i := len(qq) - 1; i >= 0; i-- {
			if !yield(qq[i]) {
				break
			}
		}
	}
}

func (pol *mlfq) prioritize(q *queue.RoundRobin) {
	for p := range drain(q) {
		pol.addToQueue(maxPriority, p)
	}
}

func drain(q *queue.RoundRobin) iter.Seq[Process] {
	return func(yield func(Process) bool) {
		for q.Len() != 0 {
			x := q.Pop()
			if x == nil {
				panic(fmt.Errorf("failed to pop"))
			}
			p, ok := x.(*process)
			if !ok {
				panic(fmt.Errorf("pop returned non-process: %v", x))
			}
			if !yield(p.proc) {
				break
			}
		}
	}
}

func (pol *mlfq) next() *process {
	for _, q := range pol.queues {
		if q.Len() == 0 {
			continue
		}
		v, ok := q.NextFunc(func(x any) bool {
			p := x.(*process)
			return !p.proc.Blocked()
		})
		if !ok {
			continue
		}
		return v.(*process)
	}
	return nil
}
