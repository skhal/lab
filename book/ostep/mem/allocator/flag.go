// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/mem/allocator/internal/heap"
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

type coalesceModeFlag struct {
	v *heap.CoalesceMode
}

func newCoalesceModeFlag(v *heap.CoalesceMode) *coalesceModeFlag {
	return &coalesceModeFlag{v}
}

var coalesceModes = []heap.CoalesceMode{
	heap.CoalesceModeNoop,
	heap.CoalesceModeForward,
	heap.CoalesceModeBackward,
	heap.CoalesceModeBidirectional,
}

// Set parses the coalesce mode from the string and stores it in the flag
// variable. It returns an error if the string does not match supported modes.
func (fl *coalesceModeFlag) Set(s string) error {
	for _, mode := range coalesceModes {
		if mode.String() == s {
			*fl.v = mode
			return nil
		}
	}
	return fmt.Errorf("unsupported mode - %s", s)
}

// String returns a string representation of the flag value.
func (fl *coalesceModeFlag) String() string {
	if fl.v == nil {
		return ""
	}
	return fl.v.String()
}

// Validate checks that the coalesce mode is set to one of the supported
// values.
func (fl *coalesceModeFlag) Validate() error {
	if !slices.Contains(coalesceModes, *fl.v) {
		return fmt.Errorf("unsupported mode - %d", fl.v)
	}
	return nil
}

type alignmentFlag struct {
	n *int
}

func newAlignmentFlag(n *int) *alignmentFlag {
	return &alignmentFlag{n}
}

// Set parses s as a number and stores it in the flag value.
func (fl *alignmentFlag) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fl.n = n
	return nil
}

// String returns a string representation of the flag value.
func (fl *alignmentFlag) String() string {
	if fl.n == nil {
		return ""
	}
	return strconv.Itoa(*fl.n)
}

// Validate checks that the alignment is either 1 (disabled) or a multiple of
// two.
func (fl *alignmentFlag) Validate() error {
	switch {
	case *fl.n == 1:
	case *fl.n%2 == 0:
	default:
		return fmt.Errorf("alignment must be multiple of 2")
	}
	return nil
}

type allocateModeFlag struct {
	v *heap.AllocateMode
}

func newAllocateModeFlag(v *heap.AllocateMode) *allocateModeFlag {
	return &allocateModeFlag{v}
}

var allocateModes = []heap.AllocateMode{
	heap.AllocateModeFirstFit,
	heap.AllocateModeNextFit,
	heap.AllocateModeBestFit,
	heap.AllocateModeWorstFit,
}

// Set parses the allocate mode from the string and stores it in the flag
// variable. It returns an error if the string does not match supported modes.
func (fl *allocateModeFlag) Set(s string) error {
	for _, mode := range allocateModes {
		if mode.String() == s {
			*fl.v = mode
			return nil
		}
	}
	return fmt.Errorf("unsupported mode - %s", s)
}

// String returns a string representation of the flag value.
func (fl *allocateModeFlag) String() string {
	if fl.v == nil {
		return ""
	}
	return fl.v.String()
}

// Validate checks that the allocate mode is set to one of the supported
// values.
func (fl *allocateModeFlag) Validate() error {
	if !slices.Contains(allocateModes, *fl.v) {
		return fmt.Errorf("unsupported mode - %d", fl.v)
	}
	return nil
}

type operationListFlag struct {
	ops *[]string
}

func newOperationListFlag(ops *[]string) *operationListFlag {
	return &operationListFlag{ops}
}

// Set parses s as a comma-separated list of strings and sets the flag value.
func (fl *operationListFlag) Set(s string) error {
	*fl.ops = append(*fl.ops, strings.Split(s, ",")...)
	return nil
}

// String returns a string representation of the flag value.
func (fl operationListFlag) String() string {
	if fl.ops == nil {
		return ""
	}
	return strings.Join(*fl.ops, ",")
}

// Validate checks that operations match "+num" malloc or "-num" free format.
func (fl operationListFlag) Validate() error {
	for i, op := range *fl.ops {
		switch {
		case strings.HasPrefix(op, "+"):
		case strings.HasPrefix(op, "-"):
		default:
			return fmt.Errorf("operation %d has invalid prefix", i)
		}
		s := strings.TrimLeft(op, "+-")
		_, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("operation %d: %s", i, err)
		}
	}
	return nil
}
