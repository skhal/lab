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

// LicenseData is input to the license template.
type LicenseData struct {
	Year   string // Year in the license block.
	Holder string // Holder in the license block.
}

func genLicenseBlock(data LicenseData) ([]byte, error) {
	var b bytes.Buffer
	if err := licenseTmpls.ExecuteTemplate(&b, "license_bsd.txt", data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
