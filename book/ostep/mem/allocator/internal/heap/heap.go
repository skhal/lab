// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap emulates memory allocations, with API simulat to malloc(3).
package heap

import (
	"encoding/binary"
	"unsafe"
)

type header struct {
	size uint16 // max 64KB
}

// Write writes header in encoded form to the buffer.
func (h *header) Write(b []byte) {
	binary.NativeEndian.PutUint16(b, h.size)
}

// Read reads and decodes the header from the buffer.
func (h *header) Read(b []byte) {
	h.size = binary.NativeEndian.Uint16(b)
}

const headerSize = int(unsafe.Sizeof(header{}))

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
func New(base, size int) *Heap {
	h := &Heap{
		base:  base,
		size:  size,
		arena: make([]byte, size),
	}
	h.free = makeFree(h.arena)
	return h
}

// makeFree initialized a block b as a free space. It write a header and
// returns offset where data buffer begins.
func makeFree(b []byte) int {
	h := header{
		size: uint16(len(b) - headerSize),
	}
	h.Write(b)
	return headerSize
}

// WalkFreeSpace walks a function f through free blocks of address space in
// the heap.
func (h *Heap) WalkFreeSpace(f func(sz, addr int) bool) {
	var hdr header
	for i := h.free; i < h.size; i += int(hdr.size) {
		hdr.Read(h.arena[i-headerSize:])
		if !f(int(hdr.size), h.base+i) {
			break
		}
	}
}
