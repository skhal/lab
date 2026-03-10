// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/lottery/internal/job"
	"github.com/skhal/lab/go/slices"
)

const (
	jobSpecSeparator     = ":"
	jobSpecListSeparator = ","
)

// JobSpecFlag parses a single job specification flag, which is a
// colon-separated list of specification fields.
type JobSpecFlag struct {
	js *job.Spec
}

// NewJobSpecFlag creates a new job specification flag for a [job.Spec] object.
// It overwrites the value of the job spec when the flag is set.
func NewJobSpecFlag(js *job.Spec) *JobSpecFlag {
	return &JobSpecFlag{js}
}

// Set implements [flag.Value] interface.
func (fl *JobSpecFlag) Set(s string) (err error) {
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			return
		}
		err = e
	}()
	const (
		idxLength  = 0
		idxTickets = 1
	)
	tt := tokens(strings.Split(s, jobSpecSeparator))
	fl.js.Length = tt.mustAtoi(idxLength, "length")
	fl.js.Tickets = tt.mustAtoi(idxTickets, "tickets")
	return nil
}

type tokens []string

func (tt tokens) mustAtoi(idx int, name string) int {
	if idx >= len(tt) {
		return 0
	}
	n, err := strconv.Atoi(tt[idx])
	if err != nil {
		panic(fmt.Errorf("parse %s: token %d: %v", name, idx, err))
	}
	return n
}

// String implements [flag.Value] interface.
func (fl *JobSpecFlag) String() string {
	if fl.js == nil {
		return ""
	}
	nn := []int{
		fl.js.Length,
		fl.js.Tickets,
	}
	return strings.Join(slices.MapFunc(nn, strconv.Itoa), jobSpecSeparator)
}

// JobSpecListFlag parses a comma-separated list of job specifications, using
// [JobSpecFlag].
type JobSpecListFlag struct {
	jsl *[]*job.Spec
	set bool
}

// NewJobSpecListFlag creates a flag for a list of job specifications.
func NewJobSpecListFlag(jsl *[]*job.Spec) *JobSpecListFlag {
	return &JobSpecListFlag{jsl: jsl}
}

// Set implements [flag.Value] interface.
func (fl *JobSpecListFlag) Set(s string) error {
	for t := range strings.SplitSeq(s, jobSpecListSeparator) {
		js := new(job.Spec)
		if err := NewJobSpecFlag(js).Set(t); err != nil {
			return err
		}
		fl.add(js)
	}
	return nil
}

func (fl *JobSpecListFlag) add(js *job.Spec) {
	if !fl.set {
		fl.set = true
		*fl.jsl = []*job.Spec{js}
		return
	}
	*fl.jsl = append(*fl.jsl, js)
}

// String implements [flag.Value] interface.
func (fl *JobSpecListFlag) String() string {
	if fl.jsl == nil {
		return ""
	}
	ss := slices.MapFunc(*fl.jsl, func(js *job.Spec) string {
		return NewJobSpecFlag(js).String()
	})
	return strings.Join(ss, jobSpecListSeparator)
}
