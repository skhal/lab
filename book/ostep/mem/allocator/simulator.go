// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/skhal/lab/book/ostep/mem/allocator/internal/heap"
)

const (
	// random points (weights) assignment to skew generation toward more
	// allocations than free.
	ptsMalloc = 6
	ptsFree   = 4

	ptsTotal = ptsMalloc + ptsFree
)

const mallocMaxSize = 2 << 10

type simulator struct {
	heap *heap.Heap

	numToGen  int // number of operations to generate
	generated int // generated number of operations

	allocated []int     // allocated addresses
	lastOp    operation // last generated operation
}

func newSimulator(h *heap.Heap, num int) *simulator {
	return &simulator{
		heap:     h,
		numToGen: num,
	}
}

// Allocated returns a slice of allocated addresses.
func (sim *simulator) Allocated() []int {
	return sim.allocated
}

// Next generates [simulator.numToGen] random operations, one at a time. It
// returns true if the operation was generated, else false.
func (sim *simulator) Next() bool {
	if sim.generated == sim.numToGen {
		sim.lastOp = nil
		return false
	}
	switch {
	case len(sim.allocated) == 0:
		// no allocated addresses available for free, malloc only.
		sim.lastOp = sim.malloc()
	default:
		sim.lastOp = sim.randOp()
	}
	sim.generated++
	return true
}

func (sim *simulator) randOp() operation {
	switch n := rand.IntN(ptsTotal); {
	case n < ptsMalloc:
		return sim.malloc()
	default:
		return sim.free()
	}
}

func (sim *simulator) malloc() operation {
	sz := 1 + rand.IntN(mallocMaxSize) // +1 for at least one byte
	return mallocOperation{
		size: sz,
		runFunc: func() error {
			a, err := sim.heap.Malloc(sz)
			if err != nil {
				return err
			}
			sim.allocated = append(sim.allocated, a)
			return nil
		},
	}
}

func (sim *simulator) free() operation {
	i := rand.IntN(len(sim.allocated))
	a := sim.allocated[i]
	return freeOperation{
		addr: a,
		runFunc: func() error {
			err := sim.heap.Free(a)
			if err != nil {
				return err
			}
			switch {
			case len(sim.allocated) == 1:
				// this operation released the only allocated address, reset the
				// allocations list.
				sim.allocated = nil
			case i == len(sim.allocated)-1:
				// released last allocation, truncate the allocations list.
				sim.allocated = sim.allocated[:i]
			default:
				copy(sim.allocated[i:], sim.allocated[i+1:])
				sim.allocated = sim.allocated[:len(sim.allocated)-1]
			}
			return nil
		},
	}
}

// Op returns the last generated operation in Next().
func (sim *simulator) Op() operation {
	return sim.lastOp
}

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
