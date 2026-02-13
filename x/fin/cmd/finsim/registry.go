// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

type registry struct {
	runners map[string]*namedRunner
	order   []string
}

func newRegistry() *registry {
	return &registry{
		runners: make(map[string]*namedRunner),
	}
}

// Get retrieves a strategy runner from the registry. It returns a boolean flag
// to indicate whether a runner with a given name is available.
func (reg *registry) Get(name string) (*namedRunner, bool) {
	r, ok := reg.runners[name]
	return r, ok
}

// Len returns the number of registered runners.
func (reg *registry) Len() int {
	return len(reg.runners)
}

// Register adds a strategy runner to the registry.
func (reg *registry) Register(r *namedRunner) error {
	name := r.Name()
	if _, ok := reg.runners[name]; ok {
		return fmt.Errorf("duplicate runner %s", name)
	}
	reg.runners[name] = r
	reg.order = append(reg.order, name)
	return nil
}

// Walk applies f to every registered strategy. The callback may return false
// to stop the iteration short.
func (reg *registry) Walk(f func(*namedRunner) bool) {
	for _, n := range reg.order {
		r := reg.runners[n]
		if !f(r) {
			break
		}
	}
}
