// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap emulates memory allocations, with API simulat to malloc(3).
package heap

import (
	"errors"
	"fmt"
	"iter"
)

var (
	// ErrAddress means passed address is invalid.
	ErrAddress = errors.New("invalid address")

	// ErrNoMemory means heap has insufficient memory to accommodate malloc().
	ErrNoMemory = errors.New("insufficient memory")

	// ErrSize means passed size is invalid.
	ErrSize = errors.New("invalid size")
)

const (
	// MinSize is the minimum supported heap size. It must be able to accommodate
	// at least the header and footer of a free block.
	MinSize = headerSize + footerSize + 1 // +1 for at lest 1B

	// MaxSize is the maximum supported heap size.
	MaxSize = int(sizeMask)
)

// Heap is a continuous address space, ready for memory allocations. It
// consists of a single free block, taking the entire heap.
type Heap struct {
	base      int // base address of the heap
	size      int // heap size
	alignment int

	enc Encoder
	dec Decoder

	scan  scanner
	coal  coalescer
	alloc allocator
}

type scanner interface {
	Scan() iter.Seq2[int, Header]
}

type coalescer interface {
	Coalesce(*Header, int)
}

type allocator interface {
	Allocate(size int) (int, error)
}

// Option is a heap option.
type Option func(hp *Heap)

// WithCoalesce option sets the heap's free space coalesce mode.
func WithCoalesce(mode CoalesceMode) Option {
	return func(hp *Heap) {
		switch mode {
		case CoalesceModeNoop:
			hp.coal = &noopCoalescer{}
		case CoalesceModeForward:
			hp.coal = newForwardCoalescer(hp)
		case CoalesceModeBackward:
			hp.coal = newBackwardCoalescer(hp)
		case CoalesceModeBidirectional:
			hp.coal = &bidiCoalescer{
				fwd: newForwardCoalescer(hp),
				bwd: newBackwardCoalescer(hp),
			}
		default:
			panic(fmt.Errorf("unsupported coalesce mode - %s", mode))
		}
	}
}

// WithAllocator option set the allocator, e.g. first fit, best fit, etc.
func WithAllocator(mode AllocateMode) Option {
	return func(hp *Heap) {
		switch mode {
		case AllocateModeBestFit:
			hp.alloc = newBestFitAllocator(hp.scan, hp.enc)
		case AllocateModeFirstFit:
			hp.alloc = newFirstFitAllocator(hp.scan, hp.enc)
		default:
			panic(fmt.Errorf("unsupported allocate mode - %s", mode))
		}
	}
}

// WithAlignment sets heap alignment. The alignment must be either 1 (no
// alignment) or a multiple of 2.
func WithAlignment(align int) Option {
	switch {
	case align == 1:
	case align%2 == 0:
	default:
		panic(fmt.Errorf("use alignment that is multiple of 2"))
	}
	return func(hp *Heap) {
		hp.alignment = align
	}
}

// New creates a heap address space of size at base address, with a single
// free block. The amount of available free space is equal to size minus the
// header size, used to store meta-information.
func New(base, size int, opts ...Option) (*Heap, error) {
	if size < MinSize {
		return nil, fmt.Errorf("%w: heap size %d, min %d", ErrSize, size, MinSize)
	}
	if size > MaxSize {
		return nil, fmt.Errorf("%w: heap size %d, max %d", ErrSize, size, MaxSize)
	}
	arena := make([]byte, size)
	hp := &Heap{
		base:      base,
		size:      size,
		alignment: 1,
		enc:       Encoder(arena),
		dec:       Decoder(arena),
		scan:      newBlockScanner(Decoder(arena), size),
		coal:      &noopCoalescer{},
		alloc:     newFirstFitAllocator(newBlockScanner(Decoder(arena), size), Encoder(arena)),
	}
	hp.enc.Encode(&Header{Size: size - headerSize}, headerSize)
	for _, opt := range opts {
		opt(hp)
	}
	return hp, nil
}

// Malloc allocates memory of requested size and returns address to newly
// allocated block. It returns zero address along with non-nil error if
// allocation fails, e.g., not enough space.
func (hp *Heap) Malloc(size int) (int, error) {
	size = hp.align(size)
	if size >= hp.size-headerSize {
		return 0, fmt.Errorf("malloc(%d): %w", size, ErrNoMemory)
	}
	if size < 1 {
		return 0, fmt.Errorf("malloc(%d): %w", size, ErrSize)
	}
	a, err := hp.alloc.Allocate(size)
	if err != nil {
		return 0, fmt.Errorf("malloc(%d): %w", size, err)
	}
	return a + hp.base, nil
}

func (hp *Heap) align(n int) int {
	m := n % hp.alignment
	if m == 0 {
		return n
	}
	padding := hp.alignment - m
	return n + padding
}

// Free releases memory at previously allocated address addr. It returns an
// errof if the address is invalid or memory is not allocated.
func (hp *Heap) Free(addr int) error {
	a := addr - hp.base
	if a < headerSize {
		return fmt.Errorf("free(%d): %w", addr, ErrAddress)
	}
	if a >= hp.size {
		return fmt.Errorf("free(%d): %w", addr, ErrAddress)
	}
	if err := hp.free(a); err != nil {
		return fmt.Errorf("free(%d): %w: %s", addr, ErrAddress, err)
	}
	return nil
}

func (hp *Heap) free(a int) error {
	var h Header
	hp.dec.Decode(&h, a)
	if !h.Allocated {
		return fmt.Errorf("block is not allocated")
	}

	h.Allocated = false
	hp.enc.Encode(&h, a)

	// remove AllocatedPrev from the header of the next block
	if b := a + h.Size + headerSize; b < hp.size {
		var hb Header
		hp.dec.Decode(&hb, b)
		hb.AllocatedPrev = false
		hp.enc.Encode(&hb, b)
	}

	hp.coal.Coalesce(&h, a)

	return nil
}

// BlockFlags
type BlockFlags struct {
	Allocated     bool // true if the block is allocated
	AllocatedPrev bool // true if the previous block is allocated
}

// Walk calls f for every block in the heap. Check block flags to distinguish
// between free and allocated blocks.
func (hp *Heap) Walk(f func(sz, addr int, fl BlockFlags)) {
	for a, h := range hp.scan.Scan() {
		f(h.Size, hp.base+a, BlockFlags{
			Allocated:     h.Allocated,
			AllocatedPrev: h.AllocatedPrev,
		})
	}
}
