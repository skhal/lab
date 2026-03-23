// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap emulates memory allocations, with API simulat to malloc(3).
package heap

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// MaxSize is the maximum heap size supported by this implementation. The value
// is a derivative of the block metadata size.
const MaxSize = int(sizeMask)

const headerSize = int(unsafe.Sizeof(uint16(0)))

type blockStatus uint16

const (
	// use iota for more status bits
	statusAllocated blockStatus = 1 << (16 - 1 - iota)

	// block status should logically AND (merge) all statuses
	statusMask blockStatus = statusAllocated
	sizeMask               = ^statusMask
)

// header holds block metadata.
type header struct {
	allocated bool
	size      int
}

// Write writes header in encoded form to the buffer.
func (h *header) Write(b []byte) {
	d := uint16(h.size & int(sizeMask))
	if h.allocated {
		d &= uint16(statusAllocated)
	}
	binary.NativeEndian.PutUint16(b, d)
}

// Read reads and decodes the header from the buffer.
func (h *header) Read(b []byte) {
	d := binary.NativeEndian.Uint16(b)
	switch blockStatus(d) & statusMask {
	case statusAllocated:
		h.allocated = true
	default:
		h.allocated = false
	}
	h.size = int(d & uint16(sizeMask))
}

// Heap is a continuous address space, ready for memory allocations. It
// consists of a single free block, taking the entire heap.
type Heap struct {
	base int // base address of the heap
	size int // heap size

	arena []byte // a continuous block of memory available to heap
	free  int    // address of the first free block inside data
}

// New creates a heap address space of size at base address, with a single
// free block. The amount of available free space is equal to size minus the
// header size, used to store meta-information.
func New(base, size int) (*Heap, error) {
	if size > MaxSize {
		return nil, fmt.Errorf("unsupported heap size %d, max %d", size, MaxSize)
	}
	h := &Heap{
		base:  base,
		size:  size,
		arena: make([]byte, size),
	}
	h.free = makeFree(h.arena)
	return h, nil
}

// makeFree initialized a block b as a free space. It write a header and
// returns offset where data buffer begins.
func makeFree(b []byte) int {
	h := header{
		size: len(b) - headerSize,
	}
	h.Write(b)
	return headerSize
}

// WalkFreeSpace walks a function f through free blocks of address space in
// the heap.
func (h *Heap) WalkFreeSpace(f func(sz, addr int) bool) {
	var hdr header
	for i := h.free; i < h.size; i += hdr.size {
		hdr.Read(h.arena[i-headerSize:])
		if !f(hdr.size, h.base+i) {
			break
		}
	}
}
