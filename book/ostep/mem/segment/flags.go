// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/book/ostep/mem/segment/internal/mem"
	"github.com/skhal/lab/go/slices"
)

const segmentFlagSeparator = ":"

// segmentFlag sets [mem.Segment.Base] and [mem.Segment.Bounds] fields from
// a flag. The flag format is <base>:<bounds>. Each number is in KB.
type segmentFlag struct {
	segment   *mem.Segment
	maxBounds mem.B
}

func newSegmentFlag(s *mem.Segment, maxBounds mem.B) *segmentFlag {
	return &segmentFlag{
		segment:   s,
		maxBounds: maxBounds,
	}
}

// Set parses segment flag into segment base and bounds fields. The flag is a
// colon-separate list of segment specification: <base>:<bounds>.
func (fl segmentFlag) Set(s string) (err error) {
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			panic(e)
		}
		err = e
	}()
	const (
		idxBase = iota
		idxBounds
	)
	splits := strings.SplitN(s, segmentFlagSeparator, 2)
	if len(splits) != 2 {
		return fmt.Errorf("invalid format, want base:bounds")
	}
	tt := tokens(splits)
	base := mem.B(tt.mustAtoi(idxBase, "base")) * mem.KB
	bounds := mem.B(tt.mustAtoi(idxBounds, "bounds")) * mem.KB
	fl.segment.Base = base
	fl.segment.Bounds = bounds
	// do not change segment direction or virtual address space offset
	return nil
}

type tokens []string

func (tt tokens) mustAtoi(idx int, name string) int {
	n, err := strconv.Atoi(tt[idx])
	if err != nil {
		panic(fmt.Errorf("%s: %v", name, err))
	}
	return n
}

// String returns a string representation of the segment flag in KB.
func (fl segmentFlag) String() string {
	if fl.segment == nil {
		return ""
	}
	nn := []mem.B{
		fl.segment.Base.KB(),
		fl.segment.Bounds.KB(),
	}
	mapfn := func(a mem.B) string {
		return strconv.Itoa(int(a))
	}
	return strings.Join(slices.MapFunc(nn, mapfn), segmentFlagSeparator)
}

// Validate checks that segmentation has valid values and returns error
// description for any violation. A valid segment has:
//   - non-negative base
//   - positive bounds
func (fl segmentFlag) Validate() error {
	if fl.segment.Base < 0 {
		return fmt.Errorf("negative base")
	}
	return newAddrBoundsFlag(&fl.segment.Bounds, 0, fl.maxBounds).Validate()
}

// addBoundsFlag sets address space bounds. It validates bounds to be
// within range [min,max].
type addBoundsFlag struct {
	bounds *mem.B
	min    mem.B
	max    mem.B
}

func newAddrBoundsFlag(b *mem.B, min, max mem.B) *addBoundsFlag {
	return &addBoundsFlag{
		bounds: b,
		min:    min,
		max:    max,
	}
}

// Set converts the bounds flag from KB to B value.
func (fl addBoundsFlag) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fl.bounds = mem.B(n) * mem.KB
	return nil
}

// String returns bounds in KB.
func (fl addBoundsFlag) String() string {
	if fl.bounds == nil {
		return ""
	}
	return strconv.Itoa(int(fl.bounds.KB()))
}

// Validate ensures that virtual address bounds are within rage of [min, max].
func (fl addBoundsFlag) Validate() error {
	if *fl.bounds < fl.min || *fl.bounds > fl.max {
		return fmt.Errorf("out of bounds [%d,%d]", fl.min, fl.max)
	}
	return nil
}
