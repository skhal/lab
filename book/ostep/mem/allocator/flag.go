// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"strconv"
)

const emptyString = ""

// ErrFlagBounds means the flag value is outside of allowed bounds.
var ErrFlagBounds = errors.New("out of bounds")

type boundedIntFlag struct {
	n   *int
	min int
	max int
}

func newBoundedIntFlag(n *int, min, max int) *boundedIntFlag {
	return &boundedIntFlag{n, min, max}
}

// Set parses flag value s and stores it in the flag variable
// [boundedIntFlag.n]. It returns an error if s string is not an integer
// number.
func (fl *boundedIntFlag) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fl.n = n
	return nil
}

// String returns a string representation of the flag variable.
func (fl *boundedIntFlag) String() string {
	if fl.n == nil {
		return emptyString
	}
	return strconv.Itoa(*fl.n)
}

// Validate checks that the flag variable is within allowed bounds of
// [min,max].
func (fl *boundedIntFlag) Validate() error {
	if *fl.n < fl.min {
		return ErrFlagBounds
	}
	if *fl.n > fl.max {
		return ErrFlagBounds
	}
	return nil
}
