// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

// State defines process state.
//
//go:generate stringer -type State -linecomment
type State int

const (
	_ State = iota

	// StateReady means the process is ready to run.
	StateReady // ready

	// StateRunning means the process is running, i.e., the last call to run the
	// process did execute the process but the process needs more CPU cycles to
	// finish.
	StateRunning // running

	// StateBlocked means the process is blocked on the IO.
	StateBlocked // blocked

	// StateZombie means the process completed.
	StateZombie // zombie
)
