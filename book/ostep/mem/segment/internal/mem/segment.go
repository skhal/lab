// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

// Direction describes segment's direction of growth.
//
//go:generate stringer -type Direction -linecomment
type Direction int

const (
	DirPositive Direction = iota // positive
	DirNegative                  // negative
)

// Segment describes a block of memory.
type Segment struct {
	Base   B // physical address of the block
	Bounds B // size of the block

	VirtBase  B         // virtual address of the block
	Direction Direction // growth direction
}
