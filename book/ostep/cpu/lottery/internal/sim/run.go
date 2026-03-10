// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
	"github.com/skhal/lab/go/slices"
)

// Run runs the simulation. It returns a list of jobs.
func Run(jsjs []*job.Spec) []*job.J {
	return slices.MapFunc(jsjs, job.New)
}
