// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type schedulerFlag struct {
	sched *sched
}

// String implements [fmt.Stringer] interface.
func (sf *schedulerFlag) String() string {
	if sf.sched == nil {
		return ""
	}
	return sf.sched.String()
}

// Set implements [flag.Value] interface.
func (sf *schedulerFlag) Set(s string) error {
	switch s {
	case "fifo":
		*sf.sched = schedFIFO
	case "sjf":
		*sf.sched = schedShortestJobFirst
	default:
		return fmt.Errorf("invalid scheduler %s", s)
	}
	return nil
}

type jobSpecFlag struct {
	specs *[]jobSpec
	set   bool
}

func newJobSpecFlag(specs *[]jobSpec) *jobSpecFlag {
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
		ss = append(ss, strconv.Itoa(spec.duration))
	}
	return strings.Join(ss, ",")
}

// Set implements [flag.Value] interface.
func (jsf *jobSpecFlag) Set(s string) error {
	defer func() { jsf.set = true }()
	if !jsf.set {
		*jsf.specs = (*jsf.specs)[:0]
	}
	for token := range strings.SplitSeq(s, ",") {
		dur, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("invalid jobs spec: %s", err)
		}
		*jsf.specs = append(*jsf.specs, jobSpec{duration: dur})
	}
	return nil
}

type jobsFlag struct {
	specs *[]jobSpec
	set   bool
}

func newJobsFlag(specs *[]jobSpec) *jobsFlag {
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
		*jf.specs = append(*jf.specs, jobSpec{duration: randomDuration})
	}
	return nil
}
