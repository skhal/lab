// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"errors"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
)

// ErrCompleted indicates that the runner completed running all cycles.
var ErrCompleted = errors.New("io runner completed")

// NewRunner creates an IO runner.
func NewRunner(dur cpu.Cycle) *Runner {
	return &Runner{
		duration: dur,
	}
}

// Runner emulates an IO that runs for a duration of CPU cycles.
type Runner struct {
	duration cpu.Cycle
	cycles   cpu.Cycle
}

// Done reports wither the [Runner] has completed IO cycles.
func (r *Runner) Done() bool {
	return r.cycles == r.duration
}

// Run accounts for a single IO cycle. It returns an error if the [Runner] is
// in [Runner.Done] state.
func (r *Runner) Run() error {
	if r.Done() {
		return ErrCompleted
	}
	r.cycles++
	return nil
}
