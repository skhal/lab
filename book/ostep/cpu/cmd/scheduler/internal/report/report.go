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
	efs   embed.FS
	tmpls = template.Must(template.New("templates").ParseFS(efs, "txt/*.txt"))
)

// Data is the report input data.
type Data struct {
	// Policy is the active scheduler's policy.
	Policy scheduler.Policy

	// JobSpec are job specifications
	JobSpecs []job.Spec

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
