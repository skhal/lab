// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scheduler

import "fmt"

const ioCost = 3 // IO takes this many cycles

const maxCycles = 100

// block identifies blocked process and what clock cycle it got blocked on.
type block struct {
	p   *Process
	clk int
}

// Scheduler manages processes
type Scheduler struct {
	processes []*Process // all processes
	ready     []*Process // processes in ready state
	running   *Process   // currently running process
	blocked   []block    // blocked processes
	done      []*Process // completed processes

	clk int // current clock cycle
}

// New creates a new Scheduler for a list of processes.
func New(pp []*Process) *Scheduler {
	s := &Scheduler{
		processes: pp,
		ready:     pp[:],
	}
	for _, p := range pp {
		p.Ready()
	}
	return s
}

// ClockCycle reports current clock cycle.
func (s *Scheduler) ClockCycle() int {
	return s.clk
}

// Step runs a single clock cycle. It returns false if the cycle was last.
func (s *Scheduler) Step() bool {
	if s.clk > maxCycles {
		panic(fmt.Errorf("reached max %d cycles", maxCycles))
	}
	if len(s.done) == len(s.processes) {
		return false
	}
	s.clk += 1
	for stop := false; !stop && s.pickRunning(); {
		s.running.Step()
		switch s.running.State() {
		case ProcessStateRunning:
			stop = true
		case ProcessStateBlocked:
			s.blockRunning()
		case ProcessStateZombie:
			s.zombieRunning()
		}
	}
	s.updateBlocked()
	return true
}

func (s *Scheduler) pickRunning() bool {
	if s.running != nil {
		return true
	}
	if ok := s.pickBlocked(); ok {
		return true
	}
	if ok := s.pickReady(); ok {
		return true
	}
	return false
}

func (s *Scheduler) pickBlocked() bool {
	if len(s.blocked) == 0 {
		return false
	}
	diff := s.clk - s.blocked[0].clk
	if diff < ioCost {
		return false
	}
	s.running = s.blocked[0].p
	s.blocked = s.blocked[1:]
	return true
}

func (s *Scheduler) pickReady() bool {
	if len(s.ready) == 0 {
		return false
	}
	s.running = s.ready[0]
	s.ready = s.ready[1:]
	return true
}

func (s *Scheduler) blockRunning() {
	s.blocked = append(s.blocked, block{p: s.running, clk: s.clk})
	s.running = nil
}

func (s *Scheduler) zombieRunning() {
	s.done = append(s.done, s.running)
	s.running = nil
}

func (s *Scheduler) updateBlocked() {
	for _, b := range s.blocked {
		if diff := s.clk - s.blocked[0].clk; diff < ioCost {
			break
		}
		if b.p.State() == ProcessStateBlocked {
			b.p.Step()
		}
	}
}
