// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/report"
	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/sim"
)

var defaultJobSpec = []*job.Spec{
	{Length: 10, Tickets: 75},
	{Length: 10, Tickets: 25},
}

// Run executes the command.
func Run() error {
	cmd := &command{defaultJobSpec}
	return cmd.Run()
}

type command struct {
	jobSpecs []*job.Spec
}

// Run parses flags, runs simulation, and prints a report to standard output.
// It returns an error if any of the stages fail.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	jobs, cycles := sim.Run(cmd.jobSpecs)
	rd := report.Data{
		Jobs:   jobs,
		Cycles: cycles,
	}
	return report.Generate(os.Stdout, rd)
}

func (cmd *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Var(NewJobSpecListFlag(&cmd.jobSpecs), "jobs", "job specifications")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	for _, s := range cmd.jobSpecs {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("error in job spec %q %w", s, err)
		}
	}
	return nil
}
