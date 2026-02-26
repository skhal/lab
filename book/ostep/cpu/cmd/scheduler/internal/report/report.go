// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"embed"
	"io"
	"text/template"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/trace"
)

var (
	//go:embed txt
	efs     embed.FS
	funcMap = template.FuncMap{
		"AverageStats": averageStats,
	}
	tmpls = template.Must(template.New("templates").Funcs(funcMap).ParseFS(efs, "txt/*.txt"))
)

// Data is the report input data.
type Data struct {
	// Policy is the active scheduler's policy.
	Policy scheduler.Policy

	// Jobs is a list of jobs in the system with ID and specification.
	Jobs []job.Job

	// Sim is the simulator reference.
	Sim *sim.Simulator

	// Tracer generates a trace from simulator. The report skips trace if the
	// tracer is nil.
	Tracer *trace.Tracer
}

// Generate creates a report and writes it to the writer.
func Generate(w io.Writer, d Data) error {
	return tmpls.ExecuteTemplate(w, "report.txt", d)
}

type avgStat struct {
	job.Stats
	count int
}

func (stat *avgStat) add(s job.Stats) {
	stat.Response += s.Response
	stat.Turnaround += s.Turnaround
	stat.Wait += s.Wait
	stat.count += 1
}

func (stat *avgStat) average() job.Stats {
	s := job.Stats{
		Response:   stat.Response / stat.count,
		Turnaround: stat.Turnaround / stat.count,
		Wait:       stat.Wait / stat.count,
	}
	return s
}

func averageStats(jobs []*job.Completed) job.Stats {
	if len(jobs) == 0 {
		return job.Stats{}
	}
	stat := new(avgStat)
	for _, j := range jobs {
		stat.add(j.Stats)
	}
	return stat.average()
}
