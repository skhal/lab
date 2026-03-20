// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

// B stands for byte in 32-bit address space.
type B int32

const (
	_  B = 1 << (10 * iota)
	KB   // kilobyte
)

// KB converts byte address to kilobytes.
func (b B) KB() B {
	return b / KB
}

// Number of bits used in segment and address mask
const (
	maskSegmentLen = 2  // number of bits used in the segment mask
	maskAddressLen = 12 // number of bits used in the address mask
)

const (
	maskSegment B = (1<<maskSegmentLen - 1) << maskAddressLen
	maskAddress B = (1<<maskAddressLen - 1)
)

// MaxVirtAddress gives the size of the virtual address space based on the
// address mask.
const MaxVirtAddress B = maskAddress + 1

// Address represents a memory address, physical or virtual, in memory. It is
// a packed pair of segment number and address in bytes.
type Address B

// MakeAddress packs segment and address b into a single block.
func MakeAddress(seg int, b B) Address {
	b += B(seg<<maskAddressLen) & maskSegment
	return Address(b)
}

// Segment pulls the segment number from the address.
func (a Address) Segment() int {
	return int((B(a) & maskSegment) >> maskAddressLen)
}

// B returns the address without segment.
func (a Address) B() B {
	return B(a) & maskAddress
}
