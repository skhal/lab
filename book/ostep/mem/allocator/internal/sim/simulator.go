// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sim simulates client application working with the heap using
// malloc(3) API: malloc() and free().
package sim

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

const mallocMaxSize = 1 << 10 // 1KB

// Simulator runs malloc operations on the heap. The operation can be random:
//
//	s := sim.NewSimulator(h, num)
//
// or manually set:
//
//	s := sim.NewSimulator(h, 0, sim.WithOps([]string{"+10"})
type Simulator struct {
	heap *heap.Heap

	numToGen  int // number of operations to generate
	generated int // generated number of operations

	allocated []int     // allocated addresses
	lastOp    operation // last generated operation

	replayer *replayer
}

// Option modifies simulator configuration in some way, e.g., set manual
// operations.
type Option func(*Simulator)

// WithOps configures simulator to replay operations.
func WithOps(ops []string) Option {
	return func(sim *Simulator) {
		sim.numToGen = len(ops)
		sim.replayer = newReplayer(ops, sim.malloc, sim.free)
	}
}

// NewSimulator creates a heap simulator to generate num random operations. Use
// WithOps to override random operations with a list of manual operation in the
// form "op,op,..." where "op" is either "+N" to allocate a block of size N or
// "-N" to free N-th currently available allocation.
func NewSimulator(h *heap.Heap, num int, opts ...Option) *Simulator {
	sim := &Simulator{
		heap:     h,
		numToGen: num,
	}
	for _, opt := range opts {
		opt(sim)
	}
	return sim
}

// Allocated returns a slice of allocated addresses.
func (sim *Simulator) Allocated() []int {
	return sim.allocated
}

// Next generates [simulator.numToGen] random operations, one at a time. It
// returns true if the operation was generated, else false.
func (sim *Simulator) Next() bool {
	if sim.generated == sim.numToGen {
		sim.lastOp = nil
		return false
	}
	switch {
	case sim.replayer != nil:
		return sim.replayer.Next()
	default:
		sim.lastOp = sim.generateOperation()
	}
	sim.generated++
	return true
}

func (sim *Simulator) generateOperation() operation {
	switch {
	case len(sim.allocated) == 0:
		// no allocated addresses available for free, malloc only.
		return sim.randMalloc()
	default:
		return sim.randOp()
	}
}

func (sim *Simulator) randOp() operation {
	switch n := rand.IntN(ptsTotal); {
	case n < ptsMalloc:
		return sim.randMalloc()
	default:
		return sim.randFree()
	}
}

func (sim *Simulator) randMalloc() operation {
	sz := 1 + rand.IntN(mallocMaxSize) // +1 for at least one byte
	return sim.malloc(sz)
}

func (sim *Simulator) malloc(sz int) operation {
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

func (sim *Simulator) randFree() operation {
	i := rand.IntN(len(sim.allocated))
	return sim.free(i)
}

func (sim *Simulator) free(idx int) operation {
	a := sim.allocated[idx]
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
			case idx == len(sim.allocated)-1:
				// released last allocation, truncate the allocations list.
				sim.allocated = sim.allocated[:idx]
			default:
				copy(sim.allocated[idx:], sim.allocated[idx+1:])
				sim.allocated = sim.allocated[:len(sim.allocated)-1]
			}
			return nil
		},
	}
}

// Op returns the last generated operation in Next().
func (sim *Simulator) Op() operation {
	if sim.replayer != nil {
		return sim.replayer.Op()
	}
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
