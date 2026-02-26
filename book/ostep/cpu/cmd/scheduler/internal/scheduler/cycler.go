// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scheduler

// Cycler counts CPU cycles.
type Cycler struct {
	num int
}

// Cycle returns current cycle.
func (c *Cycler) Cycle() int {
	return c.num
}

// Next advances to the next cycle.
func (c *Cycler) Next() {
	c.num += 1
}
