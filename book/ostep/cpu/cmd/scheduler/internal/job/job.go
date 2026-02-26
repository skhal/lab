// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

// Job is a job entity in the system with unique identifier and configuration
// settings.
type Job struct {
	// ID is a unique job identifier.
	ID int

	// Spec is the job's configuration
	Spec Spec
}
