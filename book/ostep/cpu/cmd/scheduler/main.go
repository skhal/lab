// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Scheduler implements basic CPU scheduler policies: first-in-first-out,
// shortest job run.
package main

import (
	"flag"
	"fmt"
	"iter"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	minDuration = 1
	maxDuration = 10
)

func main() {
	cmd := &command{
		JobSpecs: []jobSpec{
			{duration: randomDuration},
			{duration: randomDuration},
			{duration: randomDuration},
		},
		Policy: policyFIFO,
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type jobSpec struct {
	arrival  int
	duration int
}

type command struct {
	JobSpecs []jobSpec
	Policy   policy
	Trace    bool
}

// Run executes the command.
func (c *command) Run() error {
	if err := c.parseFlags(); err != nil {
		return err
	}
	return report.Execute(os.Stdout, struct {
		Jobs   int
		Policy policy
		Sim    *simulator
		Trace  bool
	}{
		Jobs:   len(c.JobSpecs),
		Policy: c.Policy,
		Sim:    newSimulator(c.JobSpecs, newScheduler(c.Policy)),
		Trace:  c.Trace,
	})
}

func (c *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Usage = func() {
		w := fs.Output()
		bin := filepath.Base(os.Args[0])
		fmt.Fprintf(w, "usage: %s [flags]\n", bin)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	fs.Var(newJobsFlag(&c.JobSpecs), "jobs", "number of random jobs")
	fs.Var(newJobSpecFlag(&c.JobSpecs), "job-spec", fmt.Sprintf("comma separated list of job specifications [n:]m, where n is the arrival time (default to 0) and m is the duration (%d is random)", randomDuration))
	fs.Var(&policyFlag{&c.Policy}, "policy", func() string {
		var names []string
		for _, s := range []policy{policyFIFO, policyShortestJobFirst} {
			names = append(names, s.String())
		}
		return fmt.Sprintf("scheduler policy: %s", strings.Join(names, ","))
	}())
	fs.BoolVar(&c.Trace, "trace", false, "print trace")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	validate := func(fs *flag.FlagSet) error {
		seen := make(map[string]bool)
		fs.Visit(func(f *flag.Flag) { seen[f.Name] = true })
		if seen["jobs"] && seen["job-spec"] {
			return fmt.Errorf("flags jobs and job-spec are mutually exclusive")
		}
		return nil
	}
	if err := validate(fs); err != nil {
		return err
	}
	return nil
}

// Job describes a job.
type Job struct {
	cycler *cycler

	// ID is a unique job identifier.
	ID int
	// Arrival is the cycle when the job should be added to the scheduler.
	Arrival int
	// Duration is the number of cycles the job should run for.
	Duration int

	// cycles is the number of cycles the job has run for.
	cycles int

	// arrive, run, and complete are cycles when corresponding events happen.
	arrive   int
	run      int
	complete int
}

// Arrive marks the job arrived.
func (j *Job) Arrive() {
	j.arrive = j.cycler.Cycle()
}

// Complete marks the job complete.
func (j *Job) Complete() {
	j.complete = j.cycler.Cycle()
}

// Done returns true if the job has completed the cycles, else false.
func (j *Job) Done() bool {
	return j.cycles == j.Duration
}

// Run executes the job for one cycle.
func (j *Job) Run() {
	if j.cycles == 0 {
		j.run = j.cycler.Cycle()
	}
	j.cycles += 1
}

type jobStat struct {
	Response   int
	Turnaround int
	Wait       int
}

// Stat returns job statistics including the response, turnaround, and wait
// times.
func (j *Job) Stat() jobStat {
	return jobStat{
		Response:   j.run - j.arrive,
		Turnaround: j.complete - j.arrive,
		Wait:       j.run - j.arrive,
	}
}

type scheduler interface {
	Add(*Job)
	Next() (*Job, bool)
}

type simulator struct {
	cycler *cycler
	sched  scheduler

	Jobs    []*Job
	pending []*Job
}

const randomDuration = 0

func newSimulator(jobs []jobSpec, s scheduler) *simulator {
	c := new(cycler)
	jj := make([]*Job, 0, len(jobs))
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].arrival < jobs[j].arrival
	})
	for i, job := range jobs {
		dur := job.duration
		if dur == randomDuration {
			dur = minDuration + rand.IntN(maxDuration-minDuration)
		}
		j := &Job{
			cycler:   c,
			ID:       i + 1,
			Arrival:  job.arrival,
			Duration: dur,
		}
		jj = append(jj, j)
	}
	return &simulator{
		cycler:  c,
		sched:   s,
		Jobs:    jj,
		pending: append(([]*Job)(nil), jj...),
	}
}

// Stats returns average job statistics including the response, turnaround, and
// wait times.
func (s *simulator) Stats() jobStat {
	var avg jobStat
	add := func(js jobStat) {
		avg.Response += js.Response
		avg.Turnaround += js.Turnaround
		avg.Wait += js.Wait
	}
	finalize := func() {
		n := len(s.Jobs)
		avg.Response /= n
		avg.Turnaround /= n
		avg.Wait /= n
	}
	for _, j := range s.Jobs {
		add(j.Stat())
	}
	finalize()
	return avg
}

type cycle struct {
	Num int
	Job *Job
}

// Run executes the simulator.
func (s *simulator) Run() iter.Seq[cycle] {
	return func(yield func(cycle) bool) {
		for {
			s.addJobs()
			s.cycler.Next()
			job, ok := s.sched.Next()
			if !ok && len(s.pending) == 0 {
				break
			}
			c := cycle{
				Num: s.cycler.Cycle(),
				Job: job,
			}
			if !yield(c) {
				break
			}
		}
	}
}

func (s *simulator) addJobs() {
	var (
		cycle     = s.cycler.Cycle()
		lastAdded = -1
	)
	for idx, job := range s.pending {
		if n := job.Arrival; n < cycle {
			panic(fmt.Errorf("got a job with arrival %d before cycle %d", job.Arrival, cycle))
		} else if n == cycle {
			s.sched.Add(job)
			lastAdded = idx
		} else {
			break
		}
	}
	if lastAdded >= 0 {
		s.pending = s.pending[lastAdded+1:]
	}
}

type cycler struct {
	num int
}

// Cycle returns current cycle.
func (c *cycler) Cycle() int {
	return c.num
}

// Next advances to the next cycle.
func (c *cycler) Next() {
	c.num += 1
}
