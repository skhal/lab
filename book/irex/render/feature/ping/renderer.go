// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ping renders the PingFeature.
package ping

import (
	"html/template"
	"strings"

	"github.com/skhal/lab/book/irex/pb"
)

var (
	tmpl = template.Must(template.New("ping").Parse(`
<div>ping: {{.GetPing}}</div>
<div>pong: {{.GetPong}}</div>
`))
)

type renderer struct {
	msg *pb.PingFeature
}

// NewRenderer creates a PingFeature renderer.
func NewRenderer(msg *pb.PingFeature) *renderer {
	return &renderer{
		msg: msg,
	}
}

// Render renders the PingFeature.
// It returns an error if rendering fails.
func (fr *renderer) Render() (template.HTML, error) {
	var b strings.Builder
	if err := tmpl.Execute(&b, fr.msg); err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}
