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

const (
	segmentListFlagSeparator = ","
	segmentFlagSeparator     = ":"
)

type segmentListFlag struct {
	segments *[]mem.Segment
	set      bool
}

func newSegmentListFlag(v *[]mem.Segment) *segmentListFlag {
	return &segmentListFlag{
		segments: v,
	}
}

// Set parses a comma-separated list of segments. It discards default value
// if the flag is set.
func (fl *segmentListFlag) Set(s string) error {
	for i, s := range strings.Split(s, segmentListFlagSeparator) {
		var seg mem.Segment
		if err := (segmentFlag{&seg}).Set(s); err != nil {
			return fmt.Errorf("segment %d: %v", i, err)
		}
		fl.add(seg)
	}
	return nil
}

func (fl *segmentListFlag) add(seg mem.Segment) {
	if !fl.set {
		fl.set = true
		*fl.segments = []mem.Segment{seg}
		return
	}
	*fl.segments = append(*fl.segments, seg)
}

// String returns a string representations of the segment list flag, that is a
// comma separate list of segments.
func (fl segmentListFlag) String() string {
	if fl.segments == nil {
		return ""
	}
	items := make([]string, len(*fl.segments))
	for i, segm := range *fl.segments {
		items[i] = segmentFlag{&segm}.String()
	}
	return strings.Join(items, segmentListFlagSeparator)
}

// Validate check that every segment in the list flag is correct according to
// [segmentFlag.Validate].
func (fl segmentListFlag) Validate() error {
	for i, segm := range *fl.segments {
		if err := (segmentFlag{&segm}).Validate(); err != nil {
			return fmt.Errorf("segment %d: %v", i, err)
		}
	}
	return nil
}

type segmentFlag struct {
	segment *mem.Segment
}

// Set parses segment flag into segment base and bounds fields. The flag is a
// colon-separate list of segment specification: "base:bounds".
func (fl segmentFlag) Set(s string) (err error) {
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			return
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
	var seg mem.Segment
	tt := tokens(splits)
	seg.Base = mem.Address(tt.mustAtoi(idxBase, "base")) * mem.KB
	seg.Bounds = mem.Address(tt.mustAtoi(idxBounds, "bounds")) * mem.KB
	*fl.segment = seg
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
	nn := []mem.Address{
		fl.segment.Base / mem.KB,
		fl.segment.Bounds / mem.KB,
	}
	mapfn := func(a mem.Address) string {
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
	if fl.segment.Bounds < 1 {
		return fmt.Errorf("non-positive bounds")
	}
	return nil
}
