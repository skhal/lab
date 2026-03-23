// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	"encoding/binary"
	"unsafe"
)

const headerSize = int(unsafe.Sizeof(uint16(0)))

type blockStatus uint16

const (
	// use iota for more status bits
	statusAllocated blockStatus = 1 << (16 - 1 - iota)

	// block status should logically AND (merge) all statuses
	statusMask blockStatus = statusAllocated
	sizeMask               = ^statusMask
)

// Header holds block metadata.
type Header struct {
	Allocated bool // tags block as allocated if true, else free
	Size      int  // block size
}

// Marshal encodes header to bytes. The returned slice of bytes is guaranteed
// to be of headerSize length.
func (h *Header) Marshal() []byte {
	d := uint16(h.Size & int(sizeMask))
	if h.Allocated {
		d |= uint16(statusAllocated)
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
	h.Size = int(d & uint16(sizeMask))
}

// encoder is a buffer that can encode headers.
type encoder []byte

// Encode encodes the header to the buffer at a given address.
func (e encoder) Encode(h *Header, a int) {
	copy(e[a-headerSize:a], h.Marshal())
}

// decoder is a buffer that can decode headers.
type decoder []byte

// Decode decodes a header into h for a block starting at a.
func (d decoder) Decode(h *Header, a int) {
	h.Unmarshal(d[a-headerSize : a])
}
