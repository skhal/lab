// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import "fmt"

// Control manages a job. It provides an API to change job state.
type Control struct {
	*J

	cycles int
}

// Run executes the job for one cycle. It panics if the job has completed.
func (c *Control) Run() {
	if c.Done() {
		panic(fmt.Errorf("attempt to run completed job"))
	}
	c.cycles++
	if c.cycles == c.spec.Length {
		c.done = true
	}
}
