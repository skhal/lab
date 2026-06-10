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
	// Commands are path.d commands.
	Commands []PathCommand
}

// PathCommand restricts the types of path commands.
type PathCommand interface {
	command()
}

// pathCommandBase is implements the PathCommand unexported part of the
// interface. It allows only local path commands to be added to the path.
type pathCommandBase struct{}

func (pathCommandBase) command() {}

// PathMoveCommand is the MoveTo command of SVG path. It translates to "M X Y".
type PathMoveCommand struct {
	pathCommandBase

	// Point is the destination to move to.
	Point
}

// String prints the move command.
func (cmd PathMoveCommand) String() string {
	return fmt.Sprintf("M %s", cmd.Point)
}

// PathLineCommand is the LineTo command of SVG path. It translates to "L X Y".
type PathLineCommand struct {
	pathCommandBase

	// Point is the destination to draw a line to.
	Point
}

// String prints the line command.
func (cmd PathLineCommand) String() string {
	return fmt.Sprintf("L %s", cmd.Point)
}

// String dumps path as a series of commands, starting with MoveTo command.
func (p Path) String() string {
	if len(p.Commands) == 0 {
		return ""
	}
	var b strings.Builder
	for _, c := range p.Commands {
		fmt.Fprintf(&b, " %s", c)
	}
	return b.String()[1:] // [1:] removes the leading space
}
