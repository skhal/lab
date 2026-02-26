// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"io"
	"text/template"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/trace"
)

var report *template.Template

func init() {
	const tmpl = `
{{- define "job" -}}
{{.ID}} arrival: {{.Spec.Arrival}} duration: {{.Spec.Duration}}
{{- end -}}

{{- define "stats" -}}
Response: {{.Response | printf "%-3d" }} Turnaround: {{.Turnaround | printf "%-3d" }} Wait: {{.Wait | printf "%-3d" }}
{{- end -}}

{{- define "trace" -}}
{{.Start}} run {{.Job.ID}} for {{.Cycles}} {{if eq .Cycles 1}}cycle{{else}}cycles{{end}} {{if .Job.Done}}[Done]{{end}}
{{- end -}}

policy: {{.Policy}}

jobs:
{{- range .Sim.Jobs}}
  {{template "job" .}}
{{- end}}

{{- if .Tracer}}

trace:
{{- range .Tracer.Trace}}
  {{template "trace" .}}
{{- end}}
{{- else}}{{range .Sim.Run}}{{end}}
{{- end}}

stats:
{{- range .Sim.Jobs}}
  {{.ID | printf "%-2d"}} {{template "stats" .Stat}}
{{- end}}

average:
  {{" " | printf "%2s"}} {{template "stats" .Sim.Stats}}
`
	report = template.Must(template.New("report").Parse(tmpl))
}

// Data is the report input data.
type Data struct {
	// Policy is the active scheduler's policy.
	Policy scheduler.Policy

	// Sim is the simulator reference.
	Sim *sim.Simulator

	// Tracer generates a trace from simulator. The report skips trace if the
	// tracer is nil.
	Tracer *trace.Tracer
}

// Generate creates a report and writes it to the writer.
func Generate(w io.Writer, d Data) error {
	return report.Execute(w, d)
}
