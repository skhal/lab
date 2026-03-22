// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package report generates a report for the allocator simulation. A report
// includes the configuration parameters, and simulated memory allocations.
package report

import (
	"io"
	"text/template"
)

// Data is the report input.
type Data struct {
	Heap Heap // heap configuration
}

// Heap is the heap configuration.
type Heap struct {
	Base int     // address where the heap starts.
	Size int     // size of the heap.
	Free []Block // free blocks
}

// Block is a continuous address space.
type Block struct {
	Size int // block size
	Addr int // block address
}

// Generate writes a report to w using data d. It returns an error if it fails
// to generate a report.
func Generate(w io.Writer, d Data) error {
	return tmpl.Execute(w, d)
}

var tmpl = template.Must(template.New("report").Parse(`
{{- define "heap" -}}
base: {{.Base}} size: {{.Size}} {{template "free" .Free}}
{{- end -}}

{{- define "free" -}}
free[{{len .}}]
  {{- range .}} {{template "block" .}}{{end}}
{{- end -}}

{{- define "block" -}}
{{ .Size}}:{{.Addr}}
{{- end -}}

configuration:
  {{template "heap" .Heap}}
`))
