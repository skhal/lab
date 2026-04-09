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
	AllocateModeNextFit  // next-fit
	AllocateModeBestFit  // best-fit
	AllocateModeWorstFit // worst-fit
)

type noopAllocator struct {
	enc Encoder
}

func (al *noopAllocator) allocate(a int, h *Header, size int) {
	switch {
	case h.Size > size+headerSize:
		al.split(a, *h, size)
	default:
		// not enough space to split
		h.Allocated = true
		al.enc.Encode(h, a)
	}
}

func (al *noopAllocator) split(a int, h Header, size int) {
	al.allocateBlock(a, h, size)
	offset := size + headerSize
	al.freeBlock(a+offset, h.Size-offset)
}

func (al *noopAllocator) allocateBlock(a int, h Header, size int) {
	h.Allocated = true
	h.Size = size
	al.enc.Encode(&h, a)
}

func (al *noopAllocator) freeBlock(a, size int) {
	h := Header{
		AllocatedPrev: true,
		Size:          size,
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
		switch {
		case h.Allocated: // continue searching
		case h.Size < size: // not enough space
		default:
			al.allocate(a, &h, size)
			return a, nil
		}
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
		switch {
		case h.Allocated: // continue searching
		case h.Size < size: // not enough space
		case bestFit.h == nil:
			bestFit.h = &h
			bestFit.a = a
		case h.Size < bestFit.h.Size:
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

type worstFitAllocator struct {
	*noopAllocator
	s scanner
}

func newWorstFitAllocator(s scanner, enc Encoder) *worstFitAllocator {
	return &worstFitAllocator{&noopAllocator{enc}, s}
}

// Allocate allocates memory in the largest free block that can fit the
// requested size.
func (al *worstFitAllocator) Allocate(size int) (int, error) {
	var worstFit struct {
		h *Header
		a int
	}
	for a, h := range al.s.Scan() {
		switch {
		case h.Allocated: // continue searching
		case h.Size < size: // not enough space
		case worstFit.h == nil:
			worstFit.h = &h
			worstFit.a = a
		case h.Size > worstFit.h.Size:
			worstFit.h = &h
			worstFit.a = a
		}
	}
	if worstFit.h == nil {
		return 0, ErrAllocator
	}
	al.allocate(worstFit.a, worstFit.h, size)
	return worstFit.a, nil
}
