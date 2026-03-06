// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"

	goslices "github.com/skhal/lab/go/slices"
)

var (
	// ErrPolicySpec means policy spec flag has invalid value.
	ErrPolicySpec = errors.New("invalid polisy spec")

	// ErrProcSpec means a process spec has invalid value.
	ErrProcSpec = errors.New("invalid process spec")

	// ErrProcSpecList means the process spec list has invalid process spec.
	ErrProcSpecList = errors.New("invalid process specs list")
)

const (
	specFieldSeparator = ":"
	listItemSeparator  = ","
)

// PolicySpecFlag implements [flag.Value] interface to set [policy.Spec].
type PolicySpecFlag struct {
	spec *policy.Spec
}

// NewPolicySpecFlag creates a [PolicySpecFlag] for a [policy.Spec] object.
func NewPolicySpecFlag(spec *policy.Spec) *PolicySpecFlag {
	return &PolicySpecFlag{spec}
}

// Usage returns a flag description line.
func (fg PolicySpecFlag) Usage() string {
	return "policy spec as colon-separated numbers allotment:priorities:boost-cycles"
}

type tokens []string

func (tt tokens) mustParseInt(idx int, name string) int {
	if idx >= len(tt) {
		return 0
	}
	n, err := strconv.Atoi(tt[idx])
	if err != nil {
		panic(fmt.Errorf("field %d: %s: %v", idx, name, err))
	}
	return n
}

// Set implements [flag.Value] interface. It sets the policy specification from
// the flag value s.
func (fg *PolicySpecFlag) Set(s string) (err error) {
	if s == "" {
		return
	}
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			return
		}
		err = fmt.Errorf("%w: %v", ErrPolicySpec, e)
	}()
	const (
		idxAllotment   = 0
		idxPriorities  = 1
		idxBoostCycles = 2
	)
	tt := tokens(strings.Split(s, specFieldSeparator))
	fg.spec.Allotment = cpu.Cycle(tt.mustParseInt(idxAllotment, "allotment"))
	fg.spec.Priorities = tt.mustParseInt(idxPriorities, "priorities")
	fg.spec.BoostCycles = cpu.Cycle(tt.mustParseInt(idxBoostCycles, "boost cycles"))
	return nil
}

// String implements [flag.Value] interface. It prints [policy.Spec] in flag
// parsable format.
func (fg *PolicySpecFlag) String() string {
	if fg.spec == nil {
		return ""
	}
	nn := []int{
		int(fg.spec.Allotment),
		fg.spec.Priorities,
		int(fg.spec.BoostCycles),
	}
	tokens := goslices.MapFunc(trimRightZero(nn), strconv.Itoa)
	return strings.Join(tokens, specFieldSeparator)
}

// ProcSpecListFlag implements [flag.Value] for a slice of [proc.Spec]. It
// resets the value of the list if the flag is set.
type ProcSpecListFlag struct {
	specs *[]*proc.Spec
	set   bool
}

// NewProcSpecListFlag creates a [ProcSpecListFlag].
func NewProcSpecListFlag(s *[]*proc.Spec) *ProcSpecListFlag {
	return &ProcSpecListFlag{specs: s}
}

// Set implements [flag.Value] interface.
func (fg *ProcSpecListFlag) Set(s string) error {
	if s == "" {
		return nil
	}
	var specs []*proc.Spec
	for i, s := range strings.Split(s, listItemSeparator) {
		spec := new(proc.Spec)
		flag := &ProcSpecFlag{spec}
		if err := flag.Set(s); err != nil {
			return fmt.Errorf("%w: spec %d: %w", ErrProcSpecList, i, err)
		}
		specs = append(specs, spec)
	}
	if !fg.set {
		fg.set = true
		*fg.specs = nil
	}
	*fg.specs = append(*fg.specs, specs...)
	return nil
}

// String implements [flag.Value] interface. It prints [policy.Spec] in flag
// parsable format.
func (fg *ProcSpecListFlag) String() string {
	if fg.specs == nil {
		return ""
	}
	mapfn := func(s *proc.Spec) string {
		f := &ProcSpecFlag{s}
		return f.String()
	}
	return strings.Join(goslices.MapFunc(*fg.specs, mapfn), listItemSeparator)
}

// Usage returns a flag description line.
func (fg *ProcSpecListFlag) Usage() string {
	return "process spec list as comma separated processor specs, each are a colon-separated numbers arrive:cpu-cycles:io-after-cpu-cycles:io-cycles"
}

// ProcSpecFlag is a flag for [proc.Spec].
type ProcSpecFlag struct {
	spec *proc.Spec
}

// NewProcSpecFlag creates a [ProcSpecFlag].
func NewProcSpecFlag(s *proc.Spec) *ProcSpecFlag {
	return &ProcSpecFlag{s}
}

// Set implements [flag.Value] interface.
func (fg *ProcSpecFlag) Set(s string) (err error) {
	if s == "" {
		return
	}
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			return
		}
		err = fmt.Errorf("%w: %v", ErrProcSpec, e)
	}()
	const (
		idxArrive           = 0
		idxCPUCycles        = 1
		idxIOAfterCPUCycles = 2
		idxIOCycles         = 3
	)
	tt := tokens(strings.Split(s, specFieldSeparator))
	fg.spec.Arrive = cpu.Cycle(tt.mustParseInt(idxArrive, "arrive"))
	fg.spec.CPUCycles = cpu.Cycle(tt.mustParseInt(idxCPUCycles, "cpu cycles"))
	fg.spec.IOAfterCPUCycles = cpu.Cycle(tt.mustParseInt(idxIOAfterCPUCycles, "io after cpu cycles"))
	fg.spec.IOCycles = cpu.Cycle(tt.mustParseInt(idxIOCycles, "io cycles"))
	return
}

// String implements [flag.Value] interface.
func (fg *ProcSpecFlag) String() string {
	if fg.spec == nil {
		return ""
	}
	nn := []int{
		int(fg.spec.Arrive),
		int(fg.spec.CPUCycles),
		int(fg.spec.IOAfterCPUCycles),
		int(fg.spec.IOCycles),
	}
	tokens := goslices.MapFunc(trimRightZero(nn), strconv.Itoa)
	return strings.Join(tokens, specFieldSeparator)
}

func trimRightZero(nn []int) []int {
	i := len(nn)
	for i > 0 {
		if nn[i-1] != 0 {
			break
		}
		i--
	}
	if i == 0 {
		// want at least one element
		i = 1
	}
	return nn[:i]
}
