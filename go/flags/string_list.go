// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags

import (
	"strings"
)

const listSeparator = ","

// StringList is a flag of comma-separated non-empty strings. It trims spaces
// and skips empty values. The flag wraps a slice of strings, where values
// are to be put.
//
// The wrapped slice might be non-empty, which is to be used for default value.
// If the flag is present, [StringList] overwrites the default value.
type StringList struct {
	s   *[]string
	set bool
}

// NewStringList creates a [StringList] flag.
func NewStringList(s *[]string) *StringList {
	return &StringList{s: s}
}

// Set implements [flag.Value] interface.
func (sl *StringList) Set(value string) error {
	for v := range strings.SplitSeq(value, listSeparator) {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		sl.add(v)
	}
	return nil
}

func (sl *StringList) add(s string) {
	if !sl.set {
		sl.set = true
		*sl.s = nil
	}
	*sl.s = append(*sl.s, s)
}

// String implements [flag.Value] interface.
func (sl *StringList) String() string {
	if sl.s == nil {
		return ""
	}
	return strings.Join(*sl.s, listSeparator)
}
