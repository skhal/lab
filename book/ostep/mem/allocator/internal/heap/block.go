// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	"encoding/binary"
	"fmt"
	"iter"
	"unsafe"
)

const (
	headerSize = int(unsafe.Sizeof(uint16(0)))
	footerSize = int(unsafe.Sizeof(uint16(0)))
)

type blockStatus uint16

const (
	// use iota for more status bits
	statusAllocated blockStatus = 1 << (16 - 1 - iota)
	statusAllocatedPrev

	// block status should logically AND (merge) all statuses
	statusMask blockStatus = statusAllocated | statusAllocatedPrev
	sizeMask               = ^statusMask
)

// Header holds block metadata.
type Header struct {
	Allocated     bool // tags block as allocated if true, else free
	AllocatedPrev bool // tags blocks with previous one allocated
	Size          int  // block size
}

// Marshal encodes the header to bytes. The returned slice of bytes is
// guaranteed to be of headerSize length.
func (h *Header) Marshal() []byte {
	d := uint16(h.Size & int(sizeMask))
	if h.Allocated {
		d |= uint16(statusAllocated)
	}
	if h.AllocatedPrev {
		d |= uint16(statusAllocatedPrev)
	}
	b := make([]byte, headerSize)
	binary.BigEndian.PutUint16(b, d)
	return b
}

// Unmarshal decodes header from bytes. The bytes must of headerSize length.
func (h *Header) Unmarshal(b []byte) {
	d := binary.BigEndian.Uint16(b)
	st := blockStatus(d) & statusMask
	h.Allocated = st&statusAllocated == statusAllocated
	h.AllocatedPrev = st&statusAllocatedPrev == statusAllocatedPrev
	h.Size = int(d & uint16(sizeMask))
}

// Footer holds block metadata at the end of a free block.
type Footer struct {
	Size int // size of the free block
}

// Marshal encodes the footer to bytes. The returned slice is guaranteed to
// have length equal to footerSize.
func (f *Footer) Marshal() []byte {
	d := uint16(f.Size & int(sizeMask))
	b := make([]byte, footerSize)
	binary.BigEndian.PutUint16(b, d)
	return b
}

// Unmarshal decodes footer from bytes. The bytes slice must be footerSize
// long.
func (f *Footer) Unmarshal(b []byte) {
	d := binary.BigEndian.Uint16(b)
	f.Size = int(d & uint16(sizeMask))
}

// Encoder is a buffer that can encode headers.
type Encoder []byte

// Encode encodes the header to the buffer at a given address. It adds a footer
// for free blocks.
func (e Encoder) Encode(h *Header, a int) {
	copy(e[a-headerSize:a], h.Marshal())
	if h.Allocated {
		return
	}
	f := Footer{
		Size: h.Size,
	}
	a += f.Size
	copy(e[a-footerSize:a], f.Marshal())
}

// Decoder is a buffer that can decode headers.
type Decoder []byte

// Decode decodes a header into h for a block starting at a.
func (d Decoder) Decode(h *Header, a int) {
	h.Unmarshal(d[a-headerSize : a])
}

// DecodePrevFooter decodes footer for the previous block of a block addressed
// at a.
//
//		Block N-1               Block N
//	 [header|payload|footer] [header|...]
//									 ^               ^
//									 f               a
func (d Decoder) DecodePrevFooter(f *Footer, a int) {
	a -= headerSize
	f.Unmarshal(d[a-footerSize : a])
}

type blockScanner struct {
	dec Decoder
	end int
}

func newBlockScanner(dec Decoder, end int) *blockScanner {
	return &blockScanner{dec, end}
}

// Scan runs through blocks in the heap.
func (bs *blockScanner) Scan() iter.Seq2[int, Header] {
	return func(yield func(int, Header) bool) {
		var h Header
		for a := headerSize; a < bs.end; a += h.Size + headerSize {
			bs.dec.Decode(&h, a)
			if h.Size == 0 {
				panic(fmt.Sprintf("invalid header at %d: %v", a, h))
			}
			if !yield(a, h) {
				break
			}
		}
	}
}
