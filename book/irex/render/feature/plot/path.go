// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"strings"
)

// Path is the SVG path.d attribute.
type Path struct {
	// Line is the list of LineTo commands.
	Line []Point

	// Move is the MoveTo command.
	Move Point
}

// String dumps path as a series of commands, starting with MoveTo command.
func (p Path) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "M %s", p.Move)
	for _, point := range p.Line {
		fmt.Fprintf(&b, " L %s", point)
	}
	return b.String()
}
