// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"io"
	"iter"
	"text/template"
)

var (
	//go:embed static
	efs  embed.FS
	tmpl = template.Must(template.New("report").ParseFS(efs, "static/*.txt"))
)

// ReportData holds remplate input.
type ReportData struct {
	Base             Address // offset of the physical address
	Bounds           int     // allowed size of virt address space in translations
	VirtAddressSpace int     // virtual address space size

	Frames iter.Seq[Frame] // address translations.
}

// Frame is a single virtual address translation. It holds either a physical
// address or error, the result of address translation.
type Frame struct {
	Virt Address // virtual address
	Phys Address // physica address
	Err  error   // translation error, exclusive wrt Phys
}

// Report generates a report for address translation.
func Report(w io.Writer, d ReportData) error {
	return tmpl.ExecuteTemplate(w, "memory.txt", d)
}
