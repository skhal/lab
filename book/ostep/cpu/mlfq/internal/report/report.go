// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"embed"
	"io"
	"iter"
	"text/template"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/sim"
)

var (
	//go:embed txt
	efs   embed.FS
	fnmap = template.FuncMap{
		"AvgStat": avgStat,
	}
	tmpls = template.Must(template.New("templates").Funcs(fnmap).ParseFS(efs, "txt/*.txt"))
)

// Data is the report input.
type Data struct {
	// Policy is the scheduler configuration.
	Policy policy.Spec

	// Processes are processes in the system, running, ready, completed, etc.
	Processes []Process

	// Trace is the sequence of CPU cycles.
	Trace iter.Seq[sim.Cycle]
}

// Process is the minimal interface to access process id and configuration.
type Process interface {
	// ID returns process's unique identifier.
	ID() int

	// Spec returns process's configuration.
	Spec() proc.Spec

	// Cycles return the number of completed CPU cycles.
	Cycles() cpu.Cycle

	// Stat calculate process metrics.
	Stat() proc.Stat
}

// Step generates a report with every cycle printed out.
func Step(w io.Writer, d Data) error {
	return tmpls.ExecuteTemplate(w, "step.txt", d)
}

func avgStat(pp []Process) proc.Stat {
	if len(pp) == 0 {
		return proc.Stat{}
	}
	s := proc.Stat{}
	for _, p := range pp {
		ps := p.Stat()
		s.Response += ps.Response
		s.Turnaround += ps.Turnaround
		s.Wait += ps.Wait
	}
	avg := func(n cpu.Cycle) cpu.Cycle {
		return cpu.Cycle(int(n) / len(pp))
	}
	s.Response = avg(s.Response)
	s.Turnaround = avg(s.Turnaround)
	s.Wait = avg(s.Wait)
	return s
}
