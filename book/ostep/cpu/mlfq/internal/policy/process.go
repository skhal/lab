// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

// Process is a Process interface, used by MLFQ policy.
type Process interface {
	// Arrive marks the process arrive to the system.
	Arrive()

	// Done returns true if the process completed, else false.
	Done() bool
}

type process struct {
	proc   Process
	prio   Priority
	cycles cpu.Cycle
}

func (p *process) atAllotment(allotment cpu.Cycle) bool {
	return p.cycles == allotment
}

// String implements [fmt.Stringer] interface.
func (p *process) String() string {
	return fmt.Sprintf("%s qid:%d cycles:%d", p.proc, p.prio, p.cycles)
}
