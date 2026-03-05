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
		NumQueues:   3,
		BoostCycles: 10,
	}
	defaultProcesses = []*proc.Spec{
		{CPUCycles: 10},
		{Arrive: 1, CPUCycles: 10},
		{Arrive: 2, CPUCycles: 10},
	}
)

// Run parses flags and runs simulation.
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

// Run executes the command.
func (cmd *command) Run(args []string) error {
	if err := cmd.parseFlags(args); err != nil {
		return err
	}
	pp, trace := cmd.run()
	return cmd.report(pp, trace)
}

func (cmd *command) run() ([]*proc.Process, iter.Seq[sim.Cycle]) {
	slices.SortFunc(cmd.processes, func(a, b *proc.Spec) int {
		return cmp.Compare(a.Arrive, b.Arrive)
	})
	pp := slices.Collect(goslices.MapFunc(slices.Values(cmd.processes), proc.New))
	clk := new(cpu.Clock)
	return pp, sim.Run(clk, policy.New(cmd.policy, clk), pp)
}

func (cmd *command) report(pp []*proc.Process, cc iter.Seq[sim.Cycle]) error {
	mapfn := func(p *proc.Process) report.Process {
		return report.Process(p)
	}
	data := report.Data{
		Policy:    cmd.policy,
		Processes: slices.Collect(goslices.MapFunc(slices.Values(pp), mapfn)),
		Trace:     cc,
	}
	return report.Step(os.Stdout, data)
}

func (cmd *command) parseFlags(args []string) error {
	fs := flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	// TODO(github.com/skhal/lab/issues/178): register policy flag
	// TODO(github.com/skhal/lab/issues/179): register processes flag
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	// TODO(github.com/skhal/lab/issues/178): validate policy flag
	// TODO(github.com/skhal/lab/issues/179): validate processes flag
	return nil
}
