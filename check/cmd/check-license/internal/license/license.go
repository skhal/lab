// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package license provides checker and fixer for license pre-commit check.
package license

import (
	"bytes"
	"embed"
	"text/template"
)

var (
	//go:embed data
	embedFS      embed.FS
	licenseTmpls = template.Must(template.New("licenses").ParseFS(embedFS, "data/license_*.txt"))
)

// Data is input to the license template.
type Data struct {
	Year   string // Year in the license block.
	Holder string // Holder in the license block.
}

// Generate creates a license block.
func Generate(data Data) ([]byte, error) {
	var b bytes.Buffer
	if err := licenseTmpls.ExecuteTemplate(&b, "license_bsd.txt", data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
