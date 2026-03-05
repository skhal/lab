// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

// Spec is a process configuration.
type Spec struct {
	// Arrive is the cycle when the process should appear in the system.
	Arrive int

	// CPUCycles is the number of CPU cycles the process takes to comlpete.
	CPUCycles int

	// IOAfterCPUCycles is the number of CPU cycles that must run to issue an IO.
	// Zero value means no IO is issued.
	IOAfterCPUCycles int

	// IOCycles is the number of cycles the IO runs.
	IOCycles int
}
