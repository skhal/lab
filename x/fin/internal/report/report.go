// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package report generates reports for simulated results.
package report

import (
	"embed"
	"text/template"
)

var (
	//go:embed txt
	efs   embed.FS
	tmpls = template.Must(template.New("templates").ParseFS(efs, "txt/*.txt"))
)
