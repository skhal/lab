// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
)

// ErrOutOfBounds means the address is outside allowed range.
var ErrOutOfBounds = errors.New("out of bounds address")

const minVirtAddress = 0

// Address represents a virtual or physical address.
type Address int

// String prints address in hex and dec forms. It implements [fmt.Stringer]
// interface.
func (a Address) String() string {
	return fmt.Sprintf("0x%x (%d)", int(a), a)
}

type translator struct {
	base   Address
	bounds Address
}

func newTranslator(base, bounds Address) translator {
	return translator{base, bounds}
}

// Translate interpolates virtual to physical address. It returns an error
// if the virtual address is invalid, i.g., negative or outside of bounds.
func (t translator) Translate(virt Address) (phys Address, err error) {
	if virt < minVirtAddress {
		return 0, ErrOutOfBounds
	}
	if virt >= t.bounds {
		return 0, ErrOutOfBounds
	}
	return virt + t.base, err
}
