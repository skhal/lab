// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import "fmt"

// Point is a point in Cartesian coordinates.
type Point struct {
	X int // X coordinate
	Y int // Y coordinate
}

// String prints the point, suitable for SVG path.d attribute.
func (p Point) String() string {
	return fmt.Sprintf("%d %d", p.X, p.Y)
}
