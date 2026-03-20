// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package report provides primitives to generate a simulation report.
package report

import (
	"embed"
	"io"
	"iter"
	"text/template"

	"github.com/skhal/lab/book/ostep/mem/segment/internal/mem"
)

var (
	//go:embed static
	efs  embed.FS
	tmpl = template.Must(template.New("report").ParseFS(efs, "static/*.txt"))
)

// Data contains report parameters.
type Data struct {
	VirtAddrBounds mem.B                 // size of the virtual address space
	Segments       []*mem.Segment        // virtual address segments
	Translations   iter.Seq[Translation] // virtual to physical address translations
}

// Translation is a single virtual to physical address translation. It returns
// non-nil error if the translation failed, that can be accessed through
// Error().
type Translation interface {
	Virtual() mem.Address  // virtual address
	Physical() mem.Address // physical address
	Error() error          // non-nil error if the translation failed
}

// Generate writes a report to w using data from d.
func Generate(w io.Writer, d Data) error {
	return tmpl.ExecuteTemplate(w, "report.txt", d)
}
