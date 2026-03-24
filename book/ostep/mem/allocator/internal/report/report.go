// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package report generates a report for the allocator simulation. A report
// includes the configuration parameters, and simulated memory allocations.
package report

import (
	"io"
	"iter"
	"text/template"
)

// Data is the report input.
type Data struct {
	Heap Heap                // heap configuration
	Ops  iter.Seq[Operation] // operations run by simulator
}

// Heap is the heap configuration.
type Heap struct {
	Base     int     // address where the heap starts.
	Size     int     // size of the heap.
	CoalMode string  // coalesce mode
	Free     []Block // free blocks
}

// Block is a continuous address space.
type Block struct {
	Size      int  // block size
	Addr      int  // block address
	AllocPrev bool // true if previous block is allocated
}

// Operation allocates or frees memory. It includes heap state after the
// operation runs, i.e., a list of allocated and free blocks.
type Operation struct {
	Name      string  // operation name
	Err       error   // error if any
	Addresses []int   // allocated addresses
	Allocated []Block // allocated blocks
	Free      []Block // free blocks
}

// Generate writes a report to w using data d. It returns an error if it fails
// to generate a report.
func Generate(w io.Writer, d Data) error {
	return tmpl.Execute(w, d)
}

var tmpl = template.Must(template.New("report").Parse(`
{{- define "heap" -}}
base: {{.Base}} size: {{.Size}} coalesce: {{.CoalMode}}
  {{template "free" .Free}}
{{- end -}}

{{- define "addresses" -}}
[{{len .}}] addresses
  {{- range .}} {{.}}{{end}}
{{- end -}}

{{- define "allocated" -}}
[{{len .}}] allocations
  {{- range .}} {{template "block" .}}{{end}}
{{- end -}}

{{- define "free" -}}
[{{len .}}] free blocks
  {{- range .}} {{template "block" .}}{{end}}
{{- end -}}

{{- define "operation" -}}
{{.Name}}
  {{- if .Err}} {{.Err}}
  {{- else}}
    {{template "addresses" .Addresses}}
    {{template "allocated" .Allocated}}
    {{template "free" .Free}}
  {{- end}}
{{- end -}}

{{- define "block" -}}
{{ .Size}}:{{.Addr}}[
  {{- if .AllocPrev}}P{{else}}-
  {{- end}}]
{{- end -}}

configuration:
  {{template "heap" .Heap}}

trace:
{{- range .Ops}}
  {{template "operation" .}}
{{- end}}
`))
