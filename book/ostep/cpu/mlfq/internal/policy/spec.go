// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import "github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"

// Spec is the MLFQ policy configuration.
type Spec struct {
	// Allotment is number of CPU cycles a process is allowed to run before it
	// gets de-prioritized.
	Allotment cpu.Cycle

	// Priorities is the number of priority queues in MLFQ policy.
	Priorities int

	// BoostCycles is the number of cycles a process needs to spend in the lowest
	// priority before it's priority is reset to the highest priority.
	BoostCycles cpu.Cycle
}
