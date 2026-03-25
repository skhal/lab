// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import "errors"

// ErrAllocator means the allocator failed to allocate requested memory.
var ErrAllocator = errors.New("allocator error")

// firstFitAllocator finds first block that fits the request.
type firstFitAllocator struct {
	enc encoder
	s   scanner
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
		switch {
		case h.Size > size+headerSize:
			al.split(*h, a, size)
		default:
			// not enough space to split
			h.Allocated = true
			al.enc.Encode(h, a)
		}
		return a, nil
	}
	return 0, ErrAllocator
}

func (al *firstFitAllocator) split(h Header, a int, size int) {
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
