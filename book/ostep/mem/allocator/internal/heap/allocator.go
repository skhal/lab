// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import "errors"

// ErrAllocator means the allocator failed to allocate requested memory.
var ErrAllocator = errors.New("allocator error")

// AllocateMode enumerates supported modes of free memory coalescion.
//
//go:generate stringer -type AllocateMode -linecomment
type AllocateMode int

const (
	_ AllocateMode = iota

	AllocateModeFirstFit // first-fit
	AllocateModeBestFit  // best-fit
)

type noopAllocator struct {
	enc Encoder
}

func (al *noopAllocator) allocate(a int, h *Header, size int) {
	switch {
	case h.Size > size+headerSize:
		al.split(*h, a, size)
	default:
		// not enough space to split
		h.Allocated = true
		al.enc.Encode(h, a)
	}
}

func (al *noopAllocator) split(h Header, a int, size int) {
	blockSize := h.Size

	h.Allocated = true
	h.Size = size
	al.enc.Encode(&h, a)

	// free block
	a += size + headerSize
	h = Header{
		AllocatedPrev: true,
		Size:          blockSize - size - headerSize,
	}
	al.enc.Encode(&h, a)
}

// firstFitAllocator finds first block that fits the request.
type firstFitAllocator struct {
	*noopAllocator
	s scanner
}

func newFirstFitAllocator(s scanner, enc Encoder) *firstFitAllocator {
	return &firstFitAllocator{&noopAllocator{enc}, s}
}

// Allocate finds first block that fits the requested size, splits it into
// allocated and free blocks, and returns address of the former one.
func (al *firstFitAllocator) Allocate(size int) (int, error) {
	for a, h := range al.s.Scan() {
		if h.Allocated {
			// continue searching
			continue
		}
		if h.Size < size {
			// not enough space
			continue
		}
		al.allocate(a, &h, size)
		return a, nil
	}
	return 0, ErrAllocator
}

type bestFitAllocator struct {
	*noopAllocator
	s scanner
}

func newBestFitAllocator(s scanner, enc Encoder) *bestFitAllocator {
	return &bestFitAllocator{&noopAllocator{enc}, s}
}

// Allocate allocates memory in a free block with closest size to the requested
// size.
func (al *bestFitAllocator) Allocate(size int) (int, error) {
	var bestFit struct {
		h *Header
		a int
	}
	for a, h := range al.s.Scan() {
		if h.Allocated {
			continue
		}
		if h.Size < size {
			continue
		}
		switch {
		case bestFit.h == nil:
			bestFit.h = &h
			bestFit.a = a
		case bestFit.h.Size > h.Size:
			bestFit.h = &h
			bestFit.a = a
		}
	}
	if bestFit.h == nil {
		return 0, ErrAllocator
	}
	al.allocate(bestFit.a, bestFit.h, size)
	return bestFit.a, nil
}
