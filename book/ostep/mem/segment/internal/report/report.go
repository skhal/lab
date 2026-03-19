// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package report provides primitives to generate a simulation report.
package report

import (
	"embed"
	"io"
	"text/template"

	"github.com/skhal/lab/book/ostep/mem/segment/internal/mem"
)

var (
	//go:embed static
	efs   embed.FS
	fnmap = template.FuncMap{
		"KB": func(a mem.Address) int {
			return int(a / mem.KB)
		},
	}
	tmpl = template.Must(template.New("report").Funcs(fnmap).ParseFS(efs, "static/*.txt"))
)

// Data contains report parameters.
type Data struct {
	Segments []mem.Segment // address segments.
}

// Generate writes a report to w using data from d.
func Generate(w io.Writer, d Data) error {
	return tmpl.ExecuteTemplate(w, "report.txt", d)
}
