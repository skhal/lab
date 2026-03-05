// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

// Cycle is a single CPU cycle.
type Cycle int

// Clock is a CPU clock.
type Clock struct {
	count Cycle
}

// Next advances CPU clock by one unit, i.e., cycle.
func (c *Clock) Next() {
	c.count++
}

// Cycles gives access to current CPU cycle.
func (c *Clock) Cycle() Cycle {
	return c.count
}
