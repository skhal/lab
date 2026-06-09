// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import "fmt"

// ViewBox is the SVG viewbox attribute. It defines the visible region of the
// user space, shown in the user-agent.
type ViewBox struct {
	// Point is the origin of the view box.
	Point

	// Width of the view box.
	Width int

	// Height of the view box.
	Height int
}

// String dumps the view box in format suitable for svg.viewbox attribute.
func (vb *ViewBox) String() string {
	return fmt.Sprintf("%s %d %d", vb.Point, vb.Width, vb.Height)
}
