// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Scheduler implements basic CPU schedulers: first-in-first-out, shortest
// job run, and round robin.
package main

import (
	"flag"
	"fmt"
	"iter"
	"math/rand/v2"
	"os"
	"path/filepath"
	"text/template"
)

const (
	minDuration = 1
	maxDuration = 10
)

func main() {
	cmd := &command{
		Jobs:  3,
		Sched: schedFIFO,
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//go:generate stringer -type=sched -linecomment
type sched int

const (
	_         sched = iota
	schedFIFO       // fifo
)

type schedulerFlag struct {
	sched *sched
}

// String implements [fmt.Stringer] interface.
func (sf *schedulerFlag) String() string {
	if sf.sched == nil {
		return ""
	}
	return sf.sched.String()
}

// Set implements [flag.Value] interface.
func (sf *schedulerFlag) Set(s string) error {
	switch s {
	case "fifo":
		*sf.sched = schedFIFO
	default:
		return fmt.Errorf("invalid scheduler %s", s)
	}
	return nil
}

type command struct {
	Jobs  int
	Sched sched
}

var report *template.Template

func init() {
	const tmpl = `
jobs: {{.Cmd.Jobs}}
scheduler: {{.Cmd.Sched}}

jobs:
{{- range .Sim.Jobs}}
  {{.ID}} duration: {{.Duration}}
{{- end}}

run:
{{- range .Sim.Run}}
  {{.Num | printf "%-2d"}} j{{.Job.ID}}
{{- end}}

stats:
{{- range .Sim.Jobs}}
  {{.ID | printf "%-2d"}} {{.Stat}}
{{- end}}

average:
  {{" " | printf "%2s"}} {{.Sim.Stats}}
`
	report = template.Must(template.New("report").Parse(tmpl))
}

// Run executes the command.
func (c *command) Run() error {
	if err := c.parseFlags(); err != nil {
		return err
	}
	return report.Execute(os.Stdout, struct {
		Cmd *command
		Sim *simulator
	}{
		Cmd: c,
		Sim: newSimulator(c.Jobs, newScheduler(c.Sched)),
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
	fs.IntVar(&c.Jobs, "jobs", c.Jobs, "number of jobs")
	fs.Var(&schedulerFlag{&c.Sched}, "sched", "scheduler to run")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if c.Jobs < 0 {
		return fmt.Errorf("negative jobs %d", c.Jobs)
	}
	return nil
}

// Job describes a job.
type Job struct {
	cycler *cycler

	// ID is a unique job identifier.
	ID int
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

// String implements [fmt.Stringer] interface.
func (js jobStat) String() string {
	return fmt.Sprintf("Response: %2d Turnaround %2d Wait %2d", js.Response, js.Turnaround, js.Wait)
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

func newScheduler(s sched) scheduler {
	switch s {
	case schedFIFO:
		return &fifoScheduler{}
	}
	panic(fmt.Errorf("unsupported scheduler %s", s))
}

type simulator struct {
	cycler *cycler
	sched  scheduler

	Jobs []*Job
}

func newSimulator(jobs int, s scheduler) *simulator {
	c := new(cycler)
	jj := make([]*Job, 0, jobs)
	for i := range jobs {
		j := &Job{
			cycler:   c,
			ID:       i + 1,
			Duration: minDuration + rand.IntN(maxDuration-minDuration),
		}
		jj = append(jj, j)
	}
	return &simulator{
		cycler: c,
		sched:  s,
		Jobs:   jj,
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
	for _, j := range s.Jobs {
		s.sched.Add(j)
	}
	return func(yield func(cycle) bool) {
		for j, ok := s.sched.Next(); ok; j, ok = s.sched.Next() {
			s.cycler.Next()
			c := cycle{
				Num: s.cycler.Cycle(),
				Job: j,
			}
			if !yield(c) {
				break
			}
		}
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

type fifoScheduler struct {
	cycler *cycler

	pending   []*Job
	running   *Job
	completed []*Job
}

// Add emulates job arrival to the scheduler.
func (s *fifoScheduler) Add(j *Job) {
	j.Arrive()
	s.pending = append(s.pending, j)
}

// Next returns the next job to run. The second returned parameter indicates
// whether the scheduler was able to pick up the job.
func (s *fifoScheduler) Next() (*Job, bool) {
	s.update()
	if s.running == nil {
		return nil, false
	}
	s.running.Run()
	return s.running, true
}

func (s *fifoScheduler) update() {
	if s.running != nil {
		if !s.running.Done() {
			return
		}
		s.running.Complete()
		s.completed = append(s.completed, s.running)
		s.running = nil
	}
	if len(s.pending) == 0 {
		return
	}
	s.running, s.pending = s.pending[0], s.pending[1:]
}
