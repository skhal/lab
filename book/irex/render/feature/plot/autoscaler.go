// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

// Percent is the per cent value in the range 0 to 1 inclusive.
type Percent float64

// Autoscaler adds padding to the range [xmin,xmax].
type Autoscaler struct {
	// Padding is the per cent of the range width of [xmin,max].
	Padding Percent
}

// Scale scales the range of values [xmin,xmax] to include padding.
func (scaler *Autoscaler) Scale(xmin, xmax float64) (float64, float64) {
	if scaler.Padding == 0 {
		return xmin, xmax
	}
	padding := (xmax - xmin) * float64(scaler.Padding)
	xmax += padding
	xmin = max(0, xmin-padding)
	return xmin, xmax
}
