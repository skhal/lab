// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import "fmt"

var lastid int

// J describes a job with unique identifier and specification.
type J struct {
	id   int
	spec Spec

	done bool
}

// New creates a job with unique identifier with provided specification.
func New(s *Spec) (*J, *Control) {
	lastid++
	j := &J{
		id:   lastid,
		spec: *s,
	}
	return j, &Control{J: j}
}

// Done reports whether the job completed running.
func (j *J) Done() bool {
	return j.done
}

// Spec returns job specification.
func (j *J) Spec() Spec { return j.spec }

// String implements [fmt.Stringer] interface.
func (j *J) String() string {
	return fmt.Sprintf("jid:%d %s", j.id, j.spec)
}
