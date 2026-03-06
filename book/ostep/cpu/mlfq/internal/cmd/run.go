// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"cmp"
	"flag"
	"iter"
	"os"
	"path/filepath"
	"slices"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/report"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/sim"

	goslices "github.com/skhal/lab/go/slices"
)

var (
	defaultPolicy = policy.Spec{
		Allotment:   2,
		Priorities:  3,
		BoostCycles: 10,
	}
	defaultProcesses = []*proc.Spec{
		{CPUCycles: 10},
		{Arrive: 1, CPUCycles: 8},
		{Arrive: 2, CPUCycles: 5},
	}
)

// Run initializeds the command with default flags and executes it.
func Run(args []string) error {
	cmd := &command{
		policy:    defaultPolicy,
		processes: defaultProcesses,
	}
	return cmd.Run(args)
}

type command struct {
	policy    policy.Spec
	processes []*proc.Spec
}

// Run parses the flags, runs the simulation, and generates a report.
func (cmd *command) Run(args []string) error {
	if err := cmd.parseFlags(args); err != nil {
		return err
	}
	procs, trace := cmd.run()
	return cmd.report(procs, trace)
}

func (cmd *command) run() ([]*proc.Process, iter.Seq[sim.Cycle]) {
	slices.SortFunc(cmd.processes, func(a, b *proc.Spec) int {
		return cmp.Compare(a.Arrive, b.Arrive)
	})
	clk := new(cpu.Clock)
	pol := policy.New(cmd.policy, clk)
	procs := goslices.MapFunc(cmd.processes, proc.New)
	return procs, sim.Run(clk, pol, procs)
}

func (cmd *command) report(pp []*proc.Process, cc iter.Seq[sim.Cycle]) error {
	mapfn := func(p *proc.Process) report.Process {
		return report.Process(p)
	}
	data := report.Data{
		Policy:    cmd.policy,
		Processes: goslices.MapFunc(pp, mapfn),
		Trace:     cc,
	}
	return report.Step(os.Stdout, data)
}

func (cmd *command) parseFlags(args []string) error {
	fs := flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	type valueHelper interface {
		flag.Value
		Usage() string
	}
	registerVar := func(val valueHelper, name string) {
		fs.Var(val, name, val.Usage())
	}
	registerVar(&PolicySpecFlag{&cmd.policy}, "policy")
	// TODO(github.com/skhal/lab/issues/179): register processes flag
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if err := cmd.policy.Validate(); err != nil {
		return err
	}
	// TODO(github.com/skhal/lab/issues/179): validate processes flag
	return nil
}
