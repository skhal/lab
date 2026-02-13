// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

type strategyListFlag struct {
	reg     *registry
	runners []*namedRunner

	seen map[string]bool
	set  bool
}

func newStrategyListFlag(reg *registry) *strategyListFlag {
	runners := make([]*namedRunner, 0, reg.Len())
	reg.Walk(func(r *namedRunner) bool {
		runners = append(runners, r)
		return true
	})
	return &strategyListFlag{
		reg:     reg,
		runners: runners,
		seen:    make(map[string]bool),
	}
}

// Help generates a help message for the flag.
func (f *strategyListFlag) Help() string {
	names := make([]string, 0, f.reg.Len())
	f.reg.Walk(func(r *namedRunner) bool {
		names = append(names, r.Name())
		return true
	})
	opts := strings.Join(names, "\n")
	return fmt.Sprintf("comma-separated list of strategies to run:\n%s\n", opts)
}

// Runners returns a list of registered runners.
func (f *strategyListFlag) Runners() []*namedRunner {
	return f.runners
}

// Set implements flag.Value interface.
func (f *strategyListFlag) Set(s string) error {
	var runners []*namedRunner
	for name := range strings.SplitSeq(s, ",") {
		r, ok := f.reg.Get(name)
		if !ok {
			return fmt.Errorf("unsupported strategy %s", name)
		}
		if f.seen[name] {
			return fmt.Errorf("duplicate strategy %s", name)
		}
		f.seen[name] = true
		runners = append(runners, r)
	}
	if !f.set {
		f.set = true
		f.runners = f.runners[:0]
	}
	f.runners = append(f.runners, runners...)
	return nil
}

// String implements flag.Valaue interface.
func (f *strategyListFlag) String() string {
	names := make([]string, 0, len(f.runners))
	for _, r := range f.runners {
		names = append(names, r.Name())
	}
	return strings.Join(names, ",")
}
