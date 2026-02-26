// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

// Spec is the job's configuration.
type Spec struct {
	// Arrival is the cycle when the job arrives to the scheduler.
	Arrival int
	// Duration is the number of cycles the job is expected to run.
	Duration int
}
