// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"fmt"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
)

// Cycle is a single CPU cycle.
type Cycle struct {
	Num int    // cycle number
	Job *job.J // assigned job to run
}

// String implements [fmt.Stringer] interface.
func (c Cycle) String() string {
	return fmt.Sprintf("%2d %s", c.Num, c.Job)
}
