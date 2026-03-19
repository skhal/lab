// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

// Segment describes a block of memory.
type Segment struct {
	Base   Address // physical address of the block
	Bounds Address // size of the block
}
