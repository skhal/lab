// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"fmt"
	"iter"
	"sort"

	"github.com/skhal/lab/book/ostep/cpu/cmd/sched/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/sched/internal/sched"
)

// Scheduler manages jobs.
type Scheduler interface {
	// Add introduces the job to the scheduler.
	Add(sched.Job)

	// Next retrieves the next job to run for one cycle. The second return =
	// parameter indicates whether the scheduler was able to pick up a job.
	Next() (sched.Job, bool)
}

// Simulator runs a simulation of running multiple jobs with a scheduler.
type Simulator struct {
	cycler *sched.Cycler
	sched  Scheduler

	pending   []job.Job
	completed []*job.Completed
}

// New creates a [Simulator] ready to run.
func New(jobs []job.Job, s Scheduler) *Simulator {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Spec.Arrival < jobs[j].Spec.Arrival
	})
	return &Simulator{
		cycler:  new(sched.Cycler),
		sched:   s,
		pending: jobs,
	}
}

// CompletedJobs returns the list of jobs that finished running.
func (s *Simulator) CompletedJobs() []*job.Completed {
	return s.completed
}

// Cycle emulates a CPU cycle.
type Cycle struct {
	// Num is the cycle number.
	Num int

	// Job is the job scheduled to run in this cycle.
	Job job.Job
}

// Run executes the simulator.
func (s *Simulator) Run() iter.Seq[Cycle] {
	return func(yield func(Cycle) bool) {
		for ; ; s.cycler.Next() {
			s.addJobs()
			item, ok := s.sched.Next()
			if !ok {
				if len(s.pending) == 0 {
					break
				}
				continue
			}
			lj := item.(*job.Live)
			if completed := lj.Run(); completed != nil {
				s.completed = append(s.completed, completed)
			}
			if !yield(Cycle{Num: s.cycler.Cycle(), Job: lj.Job}) {
				break
			}
		}
	}
}

func (s *Simulator) addJobs() {
	var (
		cycle     = s.cycler.Cycle()
		lastAdded = -1
	)
	for idx, j := range s.pending {
		if n := j.Spec.Arrival; n < cycle {
			panic(fmt.Errorf("got a job with arrival %d before cycle %d", j.Spec.Arrival, cycle))
		} else if n == cycle {
			lj := job.NewLive(j, s.cycler)
			s.sched.Add(lj)
			lastAdded = idx
		} else {
			break
		}
	}
	if lastAdded >= 0 {
		s.pending = s.pending[lastAdded+1:]
	}
}
