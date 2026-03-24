// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

// CoalesceMode enumerates supported modes of free memory coalescion.
//
//go:generate stringer -type CoalesceMode -linecomment
type CoalesceMode int

const (
	_ CoalesceMode = iota

	CoalesceModeNoop    // noop
	CoalesceModeForward // forward
)

type noopCoalescer struct{}

// Coalesce remains the free block untouched.
func (c *noopCoalescer) Coalesce(*Header, int) {}

type forwardCoalescer struct {
	dec    decoder
	enc    encoder
	bounds int
}

// Coalesce merges consecutive free blocks starting at a, moving forward.
func (c *forwardCoalescer) Coalesce(ha *Header, a int) {
	for {
		b := a + ha.Size + headerSize
		if b >= c.bounds {
			break
		}
		var hb Header
		c.dec.Decode(&hb, b)
		if hb.Allocated {
			break
		}
		c.coalesce(ha, &hb, a)
	}
}

func (c *forwardCoalescer) coalesce(dst, src *Header, a int) {
	dst.Size = dst.Size + headerSize + src.Size
	c.enc.Encode(dst, a)
}
