// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/cmd/process-run/internal/scheduler"
)

func main() {
	specs, simulate := mustParseFlags()
	pp := createProcesses(specs)
	dump(pp)
	if simulate {
		runSimulation(pp)
	}
}

func mustParseFlags() (specs []*scheduler.ProcessSpec, simulate bool) {
	parseFlags := func() error {
		fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
		fs.SetOutput(os.Stderr)
		fs.BoolVar(&simulate, "c", false, "run process simulation")
		fs.Var((*processSpecs)(&specs), "l", "comma-separated list of process specs instructions:probability")
		if err := fs.Parse(os.Args[1:]); err != nil {
			return err
		}
		if len(specs) == 0 {
			return errors.New("missing process spec list")
		}
		return nil
	}
	if err := parseFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

type processSpecs []*scheduler.ProcessSpec

func (f *processSpecs) Set(value string) error {
	for token := range strings.SplitSeq(value, ",") {
		spec, err := scheduler.NewProcessSpec(token)
		if err != nil {
			return err
		}
		*f = append(*f, spec)
	}
	return nil
}

func (f processSpecs) String() string {
	var ss []string
	for _, spec := range f {
		ss = append(ss, fmt.Sprintf("%s", spec))
	}
	return strings.Join(ss, ",")
}

func createProcesses(specs []*scheduler.ProcessSpec) []*scheduler.Process {
	pp := make([]*scheduler.Process, len(specs))
	for pid, spec := range specs {
		pp[pid] = scheduler.NewProcess(pid, spec)
	}
	return pp
}

func dump(processes []*scheduler.Process) {
	for _, p := range processes {
		fmt.Println("Process", p.PID())
		p.WalkInstructions(func(i scheduler.Instruction) {
			fmt.Println(" ", i)
		})
	}
}

func runSimulation(pp []*scheduler.Process) {
	r := newReporter(os.Stdout, pp)
	for s := scheduler.New(pp); s.Step(); {
		r.Report(s.ClockCycle())
	}
}

type reporter struct {
	w  io.Writer
	pp []*scheduler.Process

	printHeadline bool
}

func newReporter(w io.Writer, pp []*scheduler.Process) *reporter {
	return &reporter{
		w:             w,
		pp:            pp,
		printHeadline: true,
	}
}

func (r *reporter) Report(clk int) {
	if r.printHeadline {
		r.reportHeadline()
		r.printHeadline = false
	}
	r.reportState(clk)
}

func (r *reporter) reportHeadline() {
	fmt.Print("Clock")
	for _, p := range r.pp {
		fmt.Print("\tPID:", p.PID())
	}
	fmt.Println()
}

func (r *reporter) reportState(clk int) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%2d", clk)
	for _, p := range r.pp {
		fmt.Fprintf(&buf, "\t%s", p.State())
		last, ok := p.LastInstruction()
		if ok && p.State() != scheduler.ProcessStateZombie {
			fmt.Fprintf(&buf, ":%s", last)
		}
	}
	fmt.Println(buf.String())
}
