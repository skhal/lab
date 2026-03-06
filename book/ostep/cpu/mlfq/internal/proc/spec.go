// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"errors"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

var (
	// ErrSpecArrive means [Spec.Arrive] is negative.
	ErrSpecArrive = errors.New("invalid arrive")

	// ErrSpecCPUCycles means [Spec.CPUCycles] is non-positive.
	ErrSpecCPUCycles = errors.New("invalid cpu cycles")

	// ErrSpecIOAfterCPUCycles means [Spec.IOAfterCPUCycles] is negative.
	ErrSpecIOAfterCPUCycles = errors.New("invalid io after cpu cycles")

	// ErrSpecIOCycles means [Spec.IOAfterCPUCycles] is set but [Spec.IOCycles]
	// is non-positive.
	ErrSpecIOCycles = errors.New("invalid io cycles")
)

// Spec is a process configuration.
type Spec struct {
	// Arrive is the cycle when the process should appear in the system.
	Arrive cpu.Cycle

	// CPUCycles is the number of CPU cycles the process takes to comlpete.
	CPUCycles cpu.Cycle

	// IOAfterCPUCycles is the number of CPU cycles that must run to issue an IO.
	// Zero value means no IO is issued.
	IOAfterCPUCycles cpu.Cycle

	// IOCycles is the number of cycles the IO runs.
	IOCycles cpu.Cycle
}

// Validate ensures that [Spec] has values are:
//
// - Arrive: not negative.
// - CPUCycles: positive
// - IOAfterCPUCycles: not negative.
// - IOCycles: positive when IOAfterCPUCycles is positive.
func (s *Spec) Validate() error {
	if s.Arrive < 0 {
		return ErrSpecArrive
	}
	if s.CPUCycles < 1 {
		return ErrSpecCPUCycles
	}
	if s.IOAfterCPUCycles < 0 {
		return ErrSpecIOAfterCPUCycles
	}
	if s.IOAfterCPUCycles > 0 {
		if s.IOCycles < 1 {
			return ErrSpecIOCycles
		}
	}
	return nil
}
