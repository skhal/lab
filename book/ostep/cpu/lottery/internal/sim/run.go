// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"iter"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
)

// Run runs the simulation. It returns a list of jobs.
func Run(specs []*job.Spec) ([]*job.J, iter.Seq[Cycle]) {
	jj, cc := genJobs(specs)
	return jj, newDriver(cc).Drive()
}

func genJobs(specs []*job.Spec) ([]*job.J, []*job.Control) {
	jj := make([]*job.J, 0, len(specs))
	cc := make([]*job.Control, 0, len(specs))
	for _, js := range specs {
		j, c := job.New(js)
		jj = append(jj, j)
		cc = append(cc, c)
	}
	return jj, cc
}
