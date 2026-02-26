// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import "github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"

// Job describes a job.
type Job struct {
	// ID is a unique job identifier.
	ID int

	// Spec is the job's configuration
	Spec Spec

	state *state
}

type state struct {
	cycler *scheduler.Cycler

	runCycles int

	// arrive, run, and complete are cycles when corresponding events happen.
	cycleArrive   int
	cycleStart    int
	cycleComplete int
}

// Stat contains job metrics.
type Stat struct {
	// Response is the time from the job arrival to the time of firstrun.
	Response int

	// Turnaround is the time from the job arrival to the time of copletion.
	Turnaround int

	// Wait is the time from the job arrival to the time of firstrun.
	Wait int
}

// Duration returns the number of cycles the job would take to complete.
func (j *Job) Duration() int {
	return j.Spec.Duration
}

// LeftCycles calculates the number of cycles left for the job to run to
// completion.
func (j *Job) LeftCycles() int {
	return j.Spec.Duration - j.state.runCycles
}

// Init initializes the job state.
func (j *Job) Init(c *scheduler.Cycler) {
	j.state = &state{
		cycler:      c,
		cycleArrive: c.Cycle(),
	}
}

// Done returns true if the job has completed the cycles, else false.
func (j *Job) Done() bool {
	return j.LeftCycles() == 0
}

// Run executes the job for one cycle.
func (j *Job) Run() {
	if j.state.runCycles == 0 {
		j.state.cycleStart = j.state.cycler.Cycle()
	}
	j.state.runCycles += 1
	if j.Done() {
		j.state.cycleComplete = j.state.cycler.Cycle()
	}
}

// Stat returns job statistics including the response, turnaround, and wait
// times.
func (j *Job) Stat() Stat {
	if j.state == nil {
		return Stat{}
	}
	return Stat{
		Response:   j.state.cycleStart - j.state.cycleArrive,
		Turnaround: j.state.cycleComplete - j.state.cycleArrive + 1,
		Wait:       j.state.cycleStart - j.state.cycleArrive,
	}
}
