// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/cmd/sched/internal/sched"
)

// Live is a job in the system that was added, possibly run, but did not
// complete yet.
type Live struct {
	// Job is the job's identification and configuration data.
	Job

	cyclesCompleted int // number of cycles completed so far

	cycleArrive int // cycle when the job arrived
	cycleStart  int // cycle when the job started

	// cycler gives access to the current cycle.
	cycler *sched.Cycler
}

// NewLive creates a new Live job. It sets the job's arrive cycle to current
// cycle from the cycler.
func NewLive(j Job, c *sched.Cycler) *Live {
	return &Live{
		Job:         j,
		cycleArrive: c.Cycle(),
		cycler:      c,
	}
}

// CyclesLeft return the number of cycles left to run. It is equal to the
// duration from the spec minus the number of completed cycles.
func (j *Live) CyclesLeft() int {
	return j.Spec.Duration - j.cyclesCompleted
}

// Duration gives access to the job's spec duration.
func (j *Live) Duration() int {
	return j.Spec.Duration
}

// Done returns true if the job has completed the cycles, else false.
func (j *Live) Done() bool {
	return j.CyclesLeft() == 0
}

// Run executes the job for one cycle.
func (j *Live) Run() *Completed {
	if j.Done() {
		panic(fmt.Errorf("trying to run a completed job: %v", j.Job))
	}
	if j.cyclesCompleted == 0 {
		j.cycleStart = j.cycler.Cycle()
	}
	j.cyclesCompleted += 1
	if !j.Done() {
		return nil
	}
	return j.complete()
}

func (j *Live) complete() *Completed {
	cycle := j.cycler.Cycle()
	stats := func() Stats {
		if j.cyclesCompleted == 0 {
			return Stats{}
		}
		return Stats{
			Response:   j.cycleStart - j.cycleArrive,
			Turnaround: cycle - j.cycleArrive + 1,
			Wait:       j.cycleStart - j.cycleArrive,
		}
	}
	cj := &Completed{
		Job:   j.Job,
		Stats: stats(),
	}
	return cj
}
