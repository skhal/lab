// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

//go:generate stringer -type=policy -linecomment
type policy int

const (
	_                                   policy = iota
	policyFIFO                                 // fifo
	policyShortestJobFirst                     // sjf
	policyShortestTimeToCompletionFirst        // stcf
)

func newScheduler(s policy) *coreScheduler {
	switch s {
	case policyFIFO:
		return newCoreScheduler(fifoPolicy)
	case policyShortestJobFirst:
		return newCoreScheduler(shortestJobFirstPolicy)
	case policyShortestTimeToCompletionFirst:
		return newCoreScheduler(shortestTimeToCompletionFirstPolicy)
	}
	panic(fmt.Errorf("unsupported policy %s", s))
}

type schedState struct {
	pending   []*Job
	running   *Job
	completed []*Job
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
func (s *coreScheduler) Add(j *Job) {
	j.Arrive()
	s.state.pending = append(s.state.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *coreScheduler) Next() (*Job, bool) {
	s.update(s.state)
	if s.state.running == nil {
		return nil, false
	}
	s.state.running.Run()
	return s.state.running, true
}

func fifoPolicy(state *schedState) {
	if state.running != nil {
		if !state.running.Done() {
			return
		}
		state.running.Complete()
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
		state.running.Complete()
		state.completed = append(state.completed, state.running)
		state.running = nil
	}
	if len(state.pending) == 0 {
		return
	}
	shortest := 0
	for i, job := range state.pending {
		if job.Spec.Duration < state.pending[shortest].Spec.Duration {
			shortest = i
		}
	}
	state.running = state.pending[shortest]
	state.pending = append(state.pending[:shortest], state.pending[shortest+1:]...)
}

func shortestTimeToCompletionFirstPolicy(st *schedState) {
	if st.running != nil {
		if st.running.Done() {
			st.running.Complete()
			st.completed = append(st.completed, st.running)
			st.running = nil
		}
	}
	if len(st.pending) == 0 {
		return
	}
	next := func() (*Job, int) {
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
