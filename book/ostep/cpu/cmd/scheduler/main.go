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

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
)

const (
	minDuration = 1
	maxDuration = 10
)

func main() {
	cmd := &command{
		JobSpecs: []JobSpec{
			{Duration: randomDuration},
			{Duration: randomDuration},
			{Duration: randomDuration},
		},
		Policy: scheduler.PolicyFIFO,
	}
	if err := cmd.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// JobSpec is the job's configuration.
type JobSpec struct {
	// Arrival is the cycle when the job arrives to the scheduler.
	Arrival int
	// Duration is the number of cycles the job is expected to run.
	Duration int
}

type command struct {
	JobSpecs []JobSpec
	Policy   scheduler.Policy
	Trace    bool
}

// Run executes the command.
func (c *command) Run(args []string) error {
	if err := c.parseFlags(args); err != nil {
		return err
	}
	sim := newSimulator(c.JobSpecs, scheduler.New(c.Policy))
	tracer := func() *Tracer {
		if !c.Trace {
			return nil
		}
		return &Tracer{sim}
	}
	return report.Execute(os.Stdout, struct {
		Jobs   int
		Policy scheduler.Policy
		Sim    *simulator
		Tracer *Tracer
	}{
		Jobs:   len(c.JobSpecs),
		Policy: c.Policy,
		Sim:    sim,
		Tracer: tracer(),
	})
}

func (c *command) parseFlags(args []string) error {
	fs := flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	fs.Usage = func() {
		w := fs.Output()
		bin := filepath.Base(args[0])
		fmt.Fprintf(w, "usage: %s [flags]\n", bin)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	fs.Var(newJobsFlag(&c.JobSpecs), "jobs", "number of random jobs")
	fs.Var(newJobSpecFlag(&c.JobSpecs), "job-spec", fmt.Sprintf("comma separated list of job specifications [n:]m, where n is the arrival time (default to 0) and m is the duration (%d is random)", randomDuration))
	fs.Var(&policyFlag{&c.Policy}, "policy", func() string {
		var names []string
		for _, s := range []scheduler.Policy{
			scheduler.PolicyFIFO,
			scheduler.PolicySJF,
			scheduler.PolicySTCF,
		} {
			names = append(names, s.String())
		}
		return fmt.Sprintf("scheduler policy: %s", strings.Join(names, ","))
	}())
	fs.BoolVar(&c.Trace, "trace", false, "print trace")
	if err := fs.Parse(args[1:]); err != nil {
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
	// ID is a unique job identifier.
	ID int

	// Spec is the job's configuration
	Spec JobSpec

	state *jobState
}

// Duration returns the number of cycles the job would take to complete.
func (j *Job) Duration() int {
	return j.Spec.Duration
}

type jobState struct {
	cycler *scheduler.Cycler

	runCycles int

	// arrive, run, and complete are cycles when corresponding events happen.
	cycleArrive   int
	cycleStart    int
	cycleComplete int
}

// LeftCycles calculates the number of cycles left for the job to run to
// completion.
func (j *Job) LeftCycles() int {
	return j.Spec.Duration - j.state.runCycles
}

// Init initializes the job state.
func (j *Job) Init(c *scheduler.Cycler) {
	j.state = &jobState{
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

type jobStat struct {
	Response   int
	Turnaround int
	Wait       int
}

// Stat returns job statistics including the response, turnaround, and wait
// times.
func (j *Job) Stat() jobStat {
	if j.state == nil {
		return jobStat{}
	}
	return jobStat{
		Response:   j.state.cycleStart - j.state.cycleArrive,
		Turnaround: j.state.cycleComplete - j.state.cycleArrive + 1,
		Wait:       j.state.cycleStart - j.state.cycleArrive,
	}
}

// Scheduler manages jobs.
type Scheduler interface {
	// Add introduces the job to the scheduler.
	Add(scheduler.Job)

	// Next retrieves the next job to run for one cycle. The second return =
	// parameter indicates whether the scheduler was able to pick up a job.
	Next() (scheduler.Job, bool)
}

type simulator struct {
	cycler *scheduler.Cycler
	sched  Scheduler

	Jobs    []*Job
	pending []*Job
}

const randomDuration = 0

func newSimulator(jobs []JobSpec, s Scheduler) *simulator {
	c := new(scheduler.Cycler)
	jj := make([]*Job, 0, len(jobs))
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Arrival < jobs[j].Arrival
	})
	for i, job := range jobs {
		dur := job.Duration
		if dur == randomDuration {
			dur = minDuration + rand.IntN(maxDuration-minDuration)
		}
		j := &Job{
			ID: i + 1,
			Spec: JobSpec{
				Arrival:  job.Arrival,
				Duration: dur,
			},
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
			job, ok := s.sched.Next()
			if !ok && len(s.pending) == 0 {
				break
			}
			j := job.(*Job)
			j.Run()
			if !yield(cycle{Num: s.cycler.Cycle(), Job: j}) {
				break
			}
			s.cycler.Next()
		}
	}
}

// Trace summarizes multiple following cycles that belong to the same job. It
// describes when the run started, how many cycles it took, and what job was
// run.
type Trace struct {
	// Start is the cycle number when the trace starts.
	Start int
	// Cycles is the number of cycles of the trace.
	Cycles int
	// Job is the running job in this trace.
	Job *Job
}

// Tracer generated traces from a sequence of cycles from the simulator.
type Tracer struct {
	sim *simulator
}

// Trace generates a stream of [Trace] data.
func (t *Tracer) Trace() iter.Seq[Trace] {
	return func(yield func(Trace) bool) {
		var trace Trace
		for cycle := range t.sim.Run() {
			if trace.Job == nil {
				trace.Job = cycle.Job
			}
			if trace.Job == cycle.Job {
				trace.Cycles += 1
				continue
			}
			if !yield(trace) {
				return
			}
			trace = Trace{
				Start:  cycle.Num,
				Cycles: 1,
				Job:    cycle.Job,
			}
		}
		if trace.Job != nil {
			yield(trace)
		}
	}
}

func (s *simulator) addJobs() {
	var (
		cycle     = s.cycler.Cycle()
		lastAdded = -1
	)
	for idx, job := range s.pending {
		if n := job.Spec.Arrival; n < cycle {
			panic(fmt.Errorf("got a job with arrival %d before cycle %d", job.Spec.Arrival, cycle))
		} else if n == cycle {
			job.Init(s.cycler)
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
