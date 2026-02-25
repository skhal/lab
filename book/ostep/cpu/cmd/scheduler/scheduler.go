// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

//go:generate stringer -type=sched -linecomment
type sched int

const (
	_                     sched = iota
	schedFIFO                   // fifo
	schedShortestJobFirst       // sjf
)

type scheduler interface {
	Add(*Job)
	Next() (*Job, bool)
}

func newScheduler(s sched) scheduler {
	switch s {
	case schedFIFO:
		return &fifoScheduler{}
	case schedShortestJobFirst:
		return &shortestJobFirstScheduler{}
	}
	panic(fmt.Errorf("unsupported scheduler %s", s))
}

type fifoScheduler struct {
	cycler *cycler

	pending   []*Job
	running   *Job
	completed []*Job
}

// Add emulates job arrival to the scheduler.
func (s *fifoScheduler) Add(j *Job) {
	j.Arrive()
	s.pending = append(s.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *fifoScheduler) Next() (*Job, bool) {
	s.update()
	if s.running == nil {
		return nil, false
	}
	s.running.Run()
	return s.running, true
}

func (s *fifoScheduler) update() {
	if s.running != nil {
		if !s.running.Done() {
			return
		}
		s.running.Complete()
		s.completed = append(s.completed, s.running)
		s.running = nil
	}
	if len(s.pending) == 0 {
		return
	}
	s.running, s.pending = s.pending[0], s.pending[1:]
}

type shortestJobFirstScheduler struct {
	cycler *cycler

	pending   []*Job
	running   *Job
	completed []*Job
}

// Add emulates job arrival to the scheduler.
func (s *shortestJobFirstScheduler) Add(j *Job) {
	j.Arrive()
	s.pending = append(s.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *shortestJobFirstScheduler) Next() (*Job, bool) {
	s.update()
	if s.running == nil {
		return nil, false
	}
	s.running.Run()
	return s.running, true
}

func (s *shortestJobFirstScheduler) update() {
	if s.running != nil {
		if !s.running.Done() {
			return
		}
		s.running.Complete()
		s.completed = append(s.completed, s.running)
		s.running = nil
	}
	if len(s.pending) == 0 {
		return
	}
	shortest := 0
	for i, job := range s.pending {
		if job.Duration < s.pending[shortest].Duration {
			shortest = i
		}
	}
	s.running = s.pending[shortest]
	s.pending = append(s.pending[:shortest], s.pending[shortest+1:]...)
}
