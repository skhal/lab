// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

// Completed is a completed job with stats.
type Completed struct {
	// Job is the [Job] meta-information.
	Job

	// Stats holds job metrics.
	Stats Stats
}

// Stats contains job metrics.
type Stats struct {
	// Response is the time from the job arrival to the time of firstrun.
	Response int

	// Turnaround is the time from the job arrival to the time of copletion.
	Turnaround int

	// Wait is the time from the job arrival to the time of firstrun.
	Wait int
}
