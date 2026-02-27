// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"iter"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
)

// Trace summarizes multiple following cycles that belong to the same job. It
// describes when the run started, how many cycles it took, and what job was
// run.
type Trace struct {
	// Start is the cycle number when the trace starts.
	Start int
	// Cycles is the number of cycles of the trace.
	Cycles int
	// Job is the running job in this trace.
	Job *job.Job
}

// Tracer generated traces from a sequence of cycles from the simulator.
type Tracer struct {
	sim *sim.Simulator
}

// NewTracer creates a tracer for simulator.
func NewTracer(s *sim.Simulator) *Tracer {
	return &Tracer{
		sim: s,
	}
}

// Trace generates a stream of [Trace] data.
func (t *Tracer) Trace() iter.Seq[Trace] {
	return func(yield func(Trace) bool) {
		var trace Trace
		for cycle := range t.sim.Run() {
			if trace.Job == nil {
				trace.Start = cycle.Num
				trace.Job = &cycle.Job
			}
			if trace.Job.ID == cycle.Job.ID {
				trace.Cycles += 1
				continue
			}
			if !yield(trace) {
				return
			}
			trace = Trace{
				Start:  cycle.Num,
				Cycles: 1,
				Job:    &cycle.Job,
			}
		}
		if trace.Job != nil {
			yield(trace)
		}
	}
}
