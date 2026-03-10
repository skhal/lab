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

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/sim"
)

var (
	//go:embed static
	efs embed.FS
	tpl = template.Must(template.New("report").ParseFS(efs, "static/*.txt"))
)

// Data holds report input parameters.
type Data struct {
	Jobs   []*job.J            // jobs in the system
	Cycles iter.Seq[sim.Cycle] // cycles trace
}

// Generate creates and writes a report to w using d input parameters. It
// returns an error if report generation fails.
func Generate(w io.Writer, d Data) error {
	return tpl.ExecuteTemplate(w, "report.txt", d)
}
