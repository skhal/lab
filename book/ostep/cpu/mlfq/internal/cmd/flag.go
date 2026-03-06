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

	goslices "github.com/skhal/lab/go/slices"
)

// ErrPolicySpec means policy spec flag has invalid value.
var ErrPolicySpec = errors.New("invalid polisy spec")

const policySpecFieldSeparator = ":"

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

// Set implements [flag.Value] interface. It sets the policy specification from
// the flag value s.
func (fg *PolicySpecFlag) Set(s string) (err error) {
	if s == "" {
		return
	}
	tokens := strings.Split(s, policySpecFieldSeparator)
	const (
		idxAllotment   = 0
		idxPriorities  = 1
		idxBoostCycles = 2
	)
	mustParseInt := func(idx int, name string) int {
		if idx >= len(tokens) {
			return 0
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			panic(fmt.Errorf("%w: field %d: %s: %v", ErrPolicySpec, idx, name, err))
		}
		return n
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
		err = e
	}()
	fg.spec.Allotment = cpu.Cycle(mustParseInt(idxAllotment, "allotment"))
	fg.spec.Priorities = mustParseInt(idxPriorities, "priorities")
	fg.spec.BoostCycles = cpu.Cycle(mustParseInt(idxBoostCycles, "boost cycles"))
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
	tokens := goslices.MapFunc(nn, strconv.Itoa)
	return strings.Join(tokens, policySpecFieldSeparator)
}
