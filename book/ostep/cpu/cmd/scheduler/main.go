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
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
)

const (
	randomDuration = 0
	minDuration    = 1
	maxDuration    = 10
)

func main() {
	cmd := &command{
		JobSpecs: []job.Spec{
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

type command struct {
	JobSpecs []job.Spec
	Policy   scheduler.Policy
	Trace    bool
}

// Run executes the command.
func (c *command) Run(args []string) error {
	if err := c.parseFlags(args); err != nil {
		return err
	}
	c.randomizeJobs()
	s := sim.New(c.JobSpecs, scheduler.New(c.Policy))
	tracer := func() *Tracer {
		if !c.Trace {
			return nil
		}
		return &Tracer{s}
	}
	return report.Execute(os.Stdout, struct {
		Jobs   int
		Policy scheduler.Policy
		Sim    *sim.Simulator
		Tracer *Tracer
	}{
		Jobs:   len(c.JobSpecs),
		Policy: c.Policy,
		Sim:    s,
		Tracer: tracer(),
	})
}

func (c *command) randomizeJobs() {
	for i, spec := range c.JobSpecs {
		if spec.Duration == randomDuration {
			c.JobSpecs[i].Duration = minDuration + rand.IntN(maxDuration-minDuration)
		}
	}
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

// Trace summarizes multiple following cycles that belong to the same job. It
// describes when the run started, how many cycles it took, and what job was
// run.
type Trace struct {
	// Start is the cycle number when the trace starts.
	Start int
	// Cycles is the number of cycles of the trace.
	Cycles int
	// Job is the running job in this trace.
	Job *job.Job
}

// Tracer generated traces from a sequence of cycles from the simulator.
type Tracer struct {
	sim *sim.Simulator
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
