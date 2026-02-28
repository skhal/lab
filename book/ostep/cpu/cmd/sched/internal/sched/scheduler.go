// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched

// NewFIFO creates a scheduler with First-in-First-out policy. It runs the jobs
// based on the arrival time uninterrupted.
func NewFIFO() *scheduler {
	return newScheduler(fifoPolicy)
}

// NewSJF creates a scheduler with Shortest-Job-First policy. It runs the
// shortest job uninterrupted.
func NewSJF() *scheduler {
	return newScheduler(shortestJobFirstPolicy)
}

// NewSTCF creates a scheduler with Shortest-Time-to-Complete-First policy.
// It picks the job with shortest left work for next cycle. STCF is preemptive
// algorithm.
func NewSTCF() *scheduler {
	return newScheduler(shortestTimeToCompletionFirstPolicy)
}

// NewRoundRobin creates a scheduler with Round-Robin policy. It rotates
// pending jobs by running a new one every cycle.
func NewRoundRobin() *scheduler {
	return newScheduler(roundRobinPolicy)
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

type scheduler struct {
	state  *state
	update func(state *state)
}

func newScheduler(f updateFunc) *scheduler {
	return &scheduler{
		state:  new(state),
		update: f,
	}
}

// Add emulates job arrival to the scheduler.
func (s *scheduler) Add(j Job) {
	s.state.pending = append(s.state.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *scheduler) Next() (Job, bool) {
	if job := s.state.running; job != nil && job.Done() {
		s.state.running = nil
	}
	s.update(s.state)
	if s.state.running == nil {
		return nil, false
	}
	return s.state.running, true
}

func fifoPolicy(state *state) {
	if state.running != nil {
		return
	}
	if len(state.pending) == 0 {
		return
	}
	state.running, state.pending = state.pending[0], state.pending[1:]
}

func shortestJobFirstPolicy(state *state) {
	if state.running != nil {
		return
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

func roundRobinPolicy(st *state) {
	if len(st.pending) == 0 {
		return
	}
	if st.running == nil {
		st.running = st.pending[0]
		st.pending = st.pending[1:]
	} else {
		tmp := st.running
		st.running = st.pending[0]
		copy(st.pending, st.pending[1:])
		st.pending[len(st.pending)-1] = tmp
	}
}
