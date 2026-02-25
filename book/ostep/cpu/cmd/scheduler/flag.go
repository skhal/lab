// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
)

type policyFlag struct {
	policy *scheduler.Policy
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
		*sf.policy = scheduler.PolicyFIFO
	case "sjf":
		*sf.policy = scheduler.PolicySJF
	case "stcf":
		*sf.policy = scheduler.PolicySTCF
	default:
		return fmt.Errorf("invalid policy %s", s)
	}
	return nil
}

type jobSpecFlag struct {
	specs *[]JobSpec
	set   bool
}

func newJobSpecFlag(specs *[]JobSpec) *jobSpecFlag {
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
		*jsf.specs = append(*jsf.specs, JobSpec{Arrival: arr, Duration: dur})
	}
	return nil
}

type jobsFlag struct {
	specs *[]JobSpec
	set   bool
}

func newJobsFlag(specs *[]JobSpec) *jobsFlag {
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
		*jf.specs = append(*jf.specs, JobSpec{
			// arrival: 0 // arrive at the same time
			Duration: randomDuration,
		})
	}
	return nil
}
