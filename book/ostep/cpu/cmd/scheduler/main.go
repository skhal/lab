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
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/report"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/trace"
)

const (
	randomDuration = 0
	minDuration    = 1
	maxDuration    = 10
)

func main() {
	cmd := &command{
		jobSpecs: []job.Spec{
			{Duration: randomDuration},
			{Duration: randomDuration},
			{Duration: randomDuration},
		},
		policy: scheduler.PolicyFIFO,
	}
	if err := cmd.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type command struct {
	jobSpecs []job.Spec
	policy   scheduler.Policy
	trace    bool
}

// Run executes the command.
func (c *command) Run(args []string) error {
	if err := c.parseFlags(args); err != nil {
		return err
	}
	jobs := newJobs(c.jobSpecs)
	s := sim.New(jobs, scheduler.New(c.policy))
	tracer := func() *trace.Tracer {
		if !c.trace {
			return nil
		}
		return trace.NewTracer(s)
	}
	return report.Generate(os.Stdout, report.Data{
		Policy: c.policy.String(),
		Jobs:   jobs,
		Sim:    s,
		Tracer: tracer(),
	})
}

func newJobs(specs []job.Spec) []job.Job {
	jobs := make([]job.Job, 0, len(specs))
	for i, spec := range specs {
		if spec.Duration == randomDuration {
			spec.Duration = minDuration + rand.IntN(maxDuration-minDuration)
		}
		jobs = append(jobs, job.Job{
			ID:   i + 1, // count from 1
			Spec: spec,
		})
	}
	return jobs
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
	fs.Var(newJobsFlag(&c.jobSpecs), "jobs", "number of random jobs")
	fs.Var(newJobSpecFlag(&c.jobSpecs), "job-spec", fmt.Sprintf("comma separated list of job specifications [n:]m, where n is the arrival time (default to 0) and m is the duration (%d is random)", randomDuration))
	fs.Var(&policyFlag{&c.policy}, "policy", func() string {
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
	fs.BoolVar(&c.trace, "trace", false, "print trace")
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
