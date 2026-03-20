// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"errors"
	"fmt"
)

// ErrSegFault means address translation failed.
var ErrSegFault = errors.New("segmentation fault")

// Translator translates addresses using multiple segments. The segment number
// is encoded in [Address].
type Translator struct {
	segments []Segment
}

// NewTranslator creates an address translator with provided segments.
func NewTranslator(s ...Segment) Translator {
	return Translator{
		segments: s,
	}
}

// Translate interpolates virtual to physical address. It uses the segment
// number from the virtual address to run the translation.
func (tr Translator) Translate(vir Address) (Address, error) {
	seg, err := tr.segment(vir)
	if err != nil {
		return 0, err
	}
	offset, err := tr.offset(vir, seg)
	if err != nil {
		return 0, err
	}
	// Physical address space has no segments
	return MakeAddress(0, seg.Base+offset), nil
}

func (tr Translator) segment(vir Address) (*Segment, error) {
	n := vir.Segment()
	if n < 0 || n >= len(tr.segments) {
		return nil, fmt.Errorf("unsupported segment %d", n)
	}
	return &tr.segments[n], nil
}

func (tr Translator) offset(vir Address, seg *Segment) (B, error) {
	var offset B
	switch addr := vir.B(); seg.Direction {
	case DirPositive:
		offset = B(addr - seg.VirtBase)
		if offset < 0 || offset >= seg.Bounds {
			return 0, ErrSegFault
		}
	case DirNegative:
		offset = B(addr - seg.VirtBase)
		// can't start at 0 because at least one byte should be read from a segment
		// with negative direction
		if offset >= 0 || offset < -seg.Bounds {
			return 0, ErrSegFault
		}
		// shift by the size of the region
		offset += seg.Bounds
	default:
		panic("unreachable")
	}
	return offset, nil
}
