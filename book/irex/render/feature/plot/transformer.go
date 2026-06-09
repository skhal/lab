// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

// Transformer transforms a cartesian coordinate x to (x-offset)*scale
type Transformer struct {
	offset float64
	scale  float64
}

// NewTransformer creates a transformer with provided offset and scale.
func NewTransformer(offset float64, scale float64) *Transformer {
	return &Transformer{offset, scale}
}

// Transform transforms x to (x - offset) * scale.
func (tr *Transformer) Transform(x float64) float64 {
	return (x - tr.offset) * tr.scale
}
