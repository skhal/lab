// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap emulates memory allocations, with API simulat to malloc(3).
package heap

import (
	"errors"
	"fmt"
)

var (
	// ErrAddress mean passed address is invalid.
	ErrAddress = errors.New("invalid address")

	// ErrNoMemory means heap has insufficient memory to accommodate malloc().
	ErrNoMemory = errors.New("insufficient memory")

	// ErrSize means passed size is invalid.
	ErrSize = errors.New("invalid size")
)

// MaxSize is the maximum supported heap size. It is a derivatie of the block
// metadata size.
const MaxSize = int(sizeMask)

// Heap is a continuous address space, ready for memory allocations. It
// consists of a single free block, taking the entire heap.
type Heap struct {
	base int // base address of the heap
	size int // heap size

	enc encoder
	dec decoder

	coal coalescer
}

type coalescer interface {
	Coalesce(*Header, int)
}

type option func(hp *Heap, arena []byte)

// WithCoalesce option sets the heap's free space coalesce mode.
func WithCoalesce(mode CoalesceMode) option {
	return func(hp *Heap, arena []byte) {
		switch mode {
		case CoalesceModeNoop:
			hp.coal = &noopCoalescer{}
		case CoalesceModeForward:
			hp.coal = &forwardCoalescer{
				dec:    hp.dec,
				enc:    hp.enc,
				bounds: hp.size,
			}
		default:
			panic(fmt.Errorf("unsupported coalesce mode - %s", mode))
		}
	}
}

// New creates a heap address space of size at base address, with a single
// free block. The amount of available free space is equal to size minus the
// header size, used to store meta-information.
func New(base, size int, opts ...option) (*Heap, error) {
	if size > MaxSize {
		return nil, fmt.Errorf("unsupported heap size %d, max %d", size, MaxSize)
	}
	arena := make([]byte, size)
	hp := &Heap{
		base: base,
		size: size,
		enc:  encoder(arena),
		dec:  decoder(arena),
		coal: &noopCoalescer{},
	}
	for _, opt := range opts {
		opt(hp, arena)
	}
	h := Header{Size: size - headerSize}
	hp.enc.Encode(&h, headerSize)
	return hp, nil
}

// Malloc allocates memory of requested size and returns address to newly
// allocated block. It returns zero address along with non-nil error if
// allocation fails, e.g., not enough space.
func (hp *Heap) Malloc(size int) (int, error) {
	if size >= hp.size-headerSize {
		return 0, fmt.Errorf("malloc(%d): %w", size, ErrNoMemory)
	}
	if size < 1 {
		return 0, fmt.Errorf("malloc(%d): %w", size, ErrSize)
	}
	a, err := hp.malloc(size)
	if err != nil {
		return 0, fmt.Errorf("malloc(%d): %w", size, err)
	}
	return a + hp.base, nil
}

func (hp *Heap) malloc(size int) (addr int, err error) {
	err = ErrNoMemory
	hp.walk(func(h *Header, a int) bool {
		if h.Allocated {
			// continue searching
			return true
		}
		if h.Size < size {
			// not enough space
			return true
		}
		hp.split(a, h.Size, size)
		addr = a
		err = nil
		return false
	})
	return
}

func (hp *Heap) split(addr int, size int, n int) {
	h := Header{Allocated: true}

	takeAll := false
	if size <= n+headerSize {
		// not enough space to store header and 1+ bytes of free space
		h.Size = size
		takeAll = true
	} else {
		h.Size = n
	}
	hp.enc.Encode(&h, addr)

	if takeAll {
		return
	}

	h = Header{
		Size: size - n - headerSize,
	}
	hp.enc.Encode(&h, addr+n+headerSize)
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

	f := Footer{Size: h.Size}
	hp.enc.EncodeFooter(&f, a)

	if b := a + h.Size + headerSize; b < hp.size {
		var hb Header
		hp.dec.Decode(&hb, b)
		hb.AllocatedPrev = true
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
	hp.walk(func(h *Header, a int) bool {
		f(h.Size, hp.base+a, BlockFlags{
			Allocated:     h.Allocated,
			AllocatedPrev: h.AllocatedPrev,
		})
		return true
	})
}

func (hp *Heap) walk(f func(h *Header, addr int) bool) {
	var h Header
	for a := headerSize; a < hp.size; a += h.Size + headerSize {
		hp.dec.Decode(&h, a)
		if h.Size == 0 {
			panic(fmt.Sprintf("invalid header at %d: %v", a, h))
		}
		if !f(&h, a) {
			break
		}
	}
}
