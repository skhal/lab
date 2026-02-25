// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scheduler

import "fmt"

// Policy enumerates available scheduling policies.
//
//go:generate stringer -type=Policy -linecomment
type Policy int

const (
	_ Policy = iota
	// PolicyFIFO runs first-in-first-out job.
	PolicyFIFO // fifo
	// PolicySJF runs the job that is shortest to finish.
	PolicySJF // sjf
	// PolicySTCF preempts currently running job to pick up the shortest to
	// complete job.
	PolicySTCF // stcf
)

// New creates a scheduler for a specified policy.
func New(s Policy) *coreScheduler {
	switch s {
	case PolicyFIFO:
		return newCoreScheduler(fifoPolicy)
	case PolicySJF:
		return newCoreScheduler(shortestJobFirstPolicy)
	case PolicySTCF:
		return newCoreScheduler(shortestTimeToCompletionFirstPolicy)
	}
	panic(fmt.Errorf("unsupported policy %s", s))
}

// Job is a unit of work, managed by scheduler.
type Job interface {
	// Done reports whether the job completed.
	Done() bool

	// Duration returns the number of cycles the job would take to complete.
	// It should be job's specification duration, not outstanding number of
	// cycles.
	Duration() int

	// LeftCycles returns the number of cycles left for the job to complete.
	LeftCycles() int
}

type schedState struct {
	pending   []Job
	running   Job
	completed []Job
}

type updateFunc func(state *schedState)

type coreScheduler struct {
	state  *schedState
	update func(state *schedState)
}

func newCoreScheduler(f updateFunc) *coreScheduler {
	return &coreScheduler{
		state:  new(schedState),
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

func fifoPolicy(state *schedState) {
	if state.running != nil {
		if !state.running.Done() {
			return
		}
		state.completed = append(state.completed, state.running)
		state.running = nil
	}
	if len(state.pending) == 0 {
		return
	}
	state.running, state.pending = state.pending[0], state.pending[1:]
}

func shortestJobFirstPolicy(state *schedState) {
	if state.running != nil {
		if !state.running.Done() {
			return
		}
		state.completed = append(state.completed, state.running)
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

func shortestTimeToCompletionFirstPolicy(st *schedState) {
	if st.running != nil {
		if st.running.Done() {
			st.completed = append(st.completed, st.running)
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
			if njob.LeftCycles() <= job.LeftCycles() {
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
