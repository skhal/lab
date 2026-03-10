// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"cmp"
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
)

type driver struct {
	clk  int
	jobs []*job.Control

	tickets int

	nextJob struct {
		job *job.Control
		idx int
	}
}

func newDriver(cc []*job.Control) *driver {
	// Sort jobs by the number of tickets to speed up linear search for the job
	// with the most tickets.
	slices.SortFunc(cc, func(a, b *job.Control) int {
		return cmp.Compare(a.Spec().Tickets, b.Spec().Tickets)
	})
	d := &driver{
		jobs: cc,
	}
	for _, c := range cc {
		d.tickets += c.Spec().Tickets
	}
	return d
}

// Drive runs the simulation and generates a sequence of CPU cycles.
func (d *driver) Drive() iter.Seq[Cycle] {
	return func(yield func(Cycle) bool) {
		for d.next() && yield(d.cycle()) {
			continue
		}
	}
}

func (d *driver) done() bool { return len(d.jobs) == 0 }

func (d *driver) next() bool {
	d.updateJobs()
	if d.done() {
		return false
	}
	d.updateNext()
	return true
}

func (d *driver) updateJobs() {
	if d.nextJob.job == nil {
		return
	}
	if d.nextJob.job.Done() {
		copy(d.jobs[d.nextJob.idx:], d.jobs[d.nextJob.idx+1:])
		d.jobs = d.jobs[:len(d.jobs)-1]
		d.tickets -= d.nextJob.job.Spec().Tickets
	}
	d.nextJob.job = nil
	d.nextJob.idx = 0
}

func (d *driver) updateNext() {
	tk := 1 + rand.IntN(d.tickets) // +1 to count tickets in range [1, d.tickets]
	for i, j := range d.jobs {
		tk -= j.Spec().Tickets
		if tk > 0 {
			continue
		}
		d.nextJob.job = j
		d.nextJob.idx = i
		return
	}
	// unreachable
	panic(fmt.Errorf("can't find job"))
}

func (d *driver) cycle() Cycle {
	d.clk++
	d.nextJob.job.Run()
	return Cycle{
		Num: d.clk,
		Job: d.nextJob.job.J,
	}
}
