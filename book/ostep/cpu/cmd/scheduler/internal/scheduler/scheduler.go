// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scheduler

// NewFIFO creates a scheduler with First-in-First-out policy.
func NewFIFO() *coreScheduler {
	return newCoreScheduler(fifoPolicy)
}

// NewSJF creates a scheduler with Shortest-Job-First policy.
func NewSJF() *coreScheduler {
	return newCoreScheduler(shortestJobFirstPolicy)
}

// NewSTCF creates a scheduler with Shortest-Time-to-Complete-First policy.
func NewSTCF() *coreScheduler {
	return newCoreScheduler(shortestTimeToCompletionFirstPolicy)
}

// Job is a unit of work, managed by scheduler.
type Job interface {
	// Done reports whether the job completed.
	Done() bool

	// Duration returns the number of cycles the job would take to complete.
	// It should be job's specification duration, not outstanding number of
	// cycles.
	Duration() int

	// CyclesLeft returns the number of cycles left for the job to complete.
	CyclesLeft() int
}

type state struct {
	pending []Job
	running Job
}

type updateFunc func(state *state)

type coreScheduler struct {
	state  *state
	update func(state *state)
}

func newCoreScheduler(f updateFunc) *coreScheduler {
	return &coreScheduler{
		state:  new(state),
		update: f,
	}
}

// Add emulates job arrival to the scheduler.
func (s *coreScheduler) Add(j Job) {
	s.state.pending = append(s.state.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *coreScheduler) Next() (Job, bool) {
	s.update(s.state)
	if s.state.running == nil {
		return nil, false
	}
	return s.state.running, true
}

func fifoPolicy(state *state) {
	if state.running != nil {
		if !state.running.Done() {
			return
		}
		state.running = nil
	}
	if len(state.pending) == 0 {
		return
	}
	state.running, state.pending = state.pending[0], state.pending[1:]
}

func shortestJobFirstPolicy(state *state) {
	if state.running != nil {
		if !state.running.Done() {
			return
		}
		state.running = nil
	}
	if len(state.pending) == 0 {
		return
	}
	shortest := 0
	for i, job := range state.pending {
		if job.Duration() < state.pending[shortest].Duration() {
			shortest = i
		}
	}
	state.running = state.pending[shortest]
	state.pending = append(state.pending[:shortest], state.pending[shortest+1:]...)
}

func shortestTimeToCompletionFirstPolicy(st *state) {
	if st.running != nil {
		if st.running.Done() {
			st.running = nil
		}
	}
	if len(st.pending) == 0 {
		return
	}
	next := func() (Job, int) {
		var (
			njob = st.running
			nidx int
		)
		for i, job := range st.pending {
			if njob == nil {
				njob = job
				nidx = i
				continue
			}
			if njob.CyclesLeft() <= job.CyclesLeft() {
				continue
			}
			njob = job
			nidx = i
		}
		return njob, nidx
	}
	switch job, i := next(); job {
	case nil:
		return
	case st.running:
		return
	default:
		if st.running == nil {
			st.running = job
			st.pending = append(st.pending[:i], st.pending[i+1:]...)
		} else {
			st.running, st.pending[i] = job, st.running
		}
	}
}
