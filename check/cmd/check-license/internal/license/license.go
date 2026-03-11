// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"embed"
	"text/template"
)

var (
	//go:embed static
	efs  embed.FS
	tmpl = template.Must(template.New("licenses").ParseFS(efs, "static/*.txt"))
)

// Data is the input to the license template.
type Data struct {
	Year   string // Year in the license block.
	Holder string // Holder in the license block.
}

// Generate creates a license block.
func Generate(data Data) ([]byte, error) {
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, "bsd.txt", data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
