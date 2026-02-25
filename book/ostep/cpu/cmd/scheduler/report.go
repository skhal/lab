// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "text/template"

var report *template.Template

func init() {
	const tmpl = `
{{- define "stats" -}}
Response: {{.Response | printf "%-3d" }} Turnaround: {{.Turnaround | printf "%-3d" }} Wait: {{.Wait | printf "%-3d" }}
{{- end -}}
jobs: {{.Jobs}}
policy: {{.Policy}}

jobs:
{{- range .Sim.Jobs}}
  {{.ID}} arrival: {{.Arrival}} duration: {{.Duration}}
{{- end}}

{{- if .Trace}}

run:
{{- range .Sim.Run}}
  {{.Num | printf "%-2d"}} j{{.Job.ID}}
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
