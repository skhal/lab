// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/cmd/sched/internal/job"
)

type policyFlag struct {
	policy *Policy
}

// String implements [fmt.Stringer] interface.
func (sf *policyFlag) String() string {
	if sf.policy == nil {
		return ""
	}
	return sf.policy.String()
}

// Set implements [flag.Value] interface.
func (sf *policyFlag) Set(s string) error {
	switch s {
	case "fifo":
		*sf.policy = PolicyFIFO
	case "sjf":
		*sf.policy = PolicySJF
	case "stcf":
		*sf.policy = PolicySTCF
	case "rr":
		*sf.policy = PolicyRR
	default:
		return fmt.Errorf("invalid policy %s", s)
	}
	return nil
}

type jobSpecFlag struct {
	specs *[]job.Spec
	set   bool
}

func newJobSpecFlag(specs *[]job.Spec) *jobSpecFlag {
	return &jobSpecFlag{
		specs: specs,
	}
}

// String implements [fmt.Stringer] interface.
func (jsf *jobSpecFlag) String() string {
	if jsf.specs == nil {
		return ""
	}
	ss := make([]string, 0, len(*jsf.specs))
	for _, spec := range *jsf.specs {
		ss = append(ss, strconv.Itoa(spec.Duration))
	}
	return strings.Join(ss, ",")
}

// Set implements [flag.Value] interface.
func (jsf *jobSpecFlag) Set(s string) error {
	defer func() { jsf.set = true }()
	const (
		// spec "n:m" stands for {arrival:n, duration:m}
		idxArrival  = 0
		idxDuration = 1
	)
	if !jsf.set {
		*jsf.specs = (*jsf.specs)[:0]
	}
	for spec := range strings.SplitSeq(s, ",") {
		fields := strings.Split(spec, ":")
		switch len(fields) {
		case 1:
			fields = append([]string{"0"}, fields...)
		case 2:
		default:
			return fmt.Errorf("invalid job spec %s: want [n:]m format", spec)
		}
		arr, err := strconv.Atoi(fields[idxArrival])
		if err != nil {
			return fmt.Errorf("invalid job spec %s: %s", spec, err)
		}
		dur, err := strconv.Atoi(fields[idxDuration])
		if err != nil {
			return fmt.Errorf("invalid job spec %s: %s", spec, err)
		}
		*jsf.specs = append(*jsf.specs, job.Spec{Arrival: arr, Duration: dur})
	}
	return nil
}

type jobsFlag struct {
	specs *[]job.Spec
	set   bool
}

func newJobsFlag(specs *[]job.Spec) *jobsFlag {
	return &jobsFlag{
		specs: specs,
	}
}

// String implements [fmt.Stringer] interface.
func (jf *jobsFlag) String() string {
	if jf.specs == nil {
		return ""
	}
	return strconv.Itoa(len(*jf.specs))
}

// Set implements [flag.Value] interface.
func (jf *jobsFlag) Set(s string) error {
	defer func() { jf.set = true }()
	if !jf.set {
		*jf.specs = (*jf.specs)[:0]
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	for range n {
		*jf.specs = append(*jf.specs, job.Spec{
			// arrival: 0 // arrive at the same time
			Duration: randomDuration,
		})
	}
	return nil
}
