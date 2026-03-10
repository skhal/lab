// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"fmt"
	"iter"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/sim"
)

// Run generates a sequence of cycle blocks. Each cycle block describes a
// sequence of following cycles that correspond to the same job.
func Run(cc iter.Seq[sim.Cycle]) iter.Seq[CycleBlock] {
	return func(yield func(CycleBlock) bool) {
		var cb CycleBlock
		for c := range cc {
			switch cb.Job {
			case nil:
				cb.Cycle.First = c.Num
				cb.Cycle.Last = c.Num
				cb.Job = c.Job
				continue
			case c.Job:
				cb.Cycle.Last = c.Num
				continue
			}
			if !yield(cb) {
				return
			}
			cb.Job = c.Job
			cb.Cycle.First = c.Num
			cb.Cycle.Last = c.Num
		}
		if cb.Job != nil {
			yield(cb)
		}
	}
}

// CycleBlock describe a sequence of cycles that belong to the same job.
type CycleBlock struct {
	// Cycle stores the first and last cycle number for the sequence.
	Cycle struct {
		First int // first cycle in the sequence
		Last  int // last cycle in the sequence
	}
	Job *job.J // job that was run in the sequence
}

// String implements [fmt.Stringer] interface.
func (cb CycleBlock) String() string {
	if d := cb.Cycle.Last - cb.Cycle.First; d > 0 {
		return fmt.Sprintf("%-2d +%d %s", cb.Cycle.First, d, cb.Job)
	}
	return fmt.Sprintf("%-5d %s", cb.Cycle.First, cb.Job)
}
