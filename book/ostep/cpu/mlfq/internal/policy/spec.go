// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"errors"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

var (
	// ErrAllotment means [Spec.Allotment] field has non-positive value.
	ErrAllotment = errors.New("invalid allotment")

	// ErrPriorities means [Spec.Priorities] field has non-positive value.
	ErrPriorities = errors.New("invalid priorities")

	// ErrBoostCycles means [Spec.BoostCycles] field has non-positive value.
	ErrBoostCycles = errors.New("invalid boost cycles")
)

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

// Validate verifies that the specification has fields with positive value. It
// returns an error for the first invalid field if any else nil.
func (s *Spec) Validate() error {
	if s.Allotment < 1 {
		return ErrAllotment
	}
	if s.Priorities < 1 {
		return ErrPriorities
	}
	if s.BoostCycles < 1 {
		return ErrBoostCycles
	}
	return nil
}
