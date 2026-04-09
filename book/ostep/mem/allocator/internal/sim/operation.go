// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import "fmt"

type operation interface {
	fmt.Stringer
	Run() error // execute the operation.
}

type mallocOperation struct {
	size int
	runFunc
}

// String returns operation name.
func (op mallocOperation) String() string {
	return fmt.Sprintf("malloc(%d)", op.size)
}

type freeOperation struct {
	addr int
	runFunc
}

// String returns operation name.
func (op freeOperation) String() string {
	return fmt.Sprintf("free(%d)", op.addr)
}

type runFunc func() error

// Run executes the operation.
func (r runFunc) Run() error {
	return r()
}
