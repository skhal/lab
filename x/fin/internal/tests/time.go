// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"testing"
	"time"
)

// NewTime testing helper creates a time with year and month.
func NewTime(t *testing.T, year int, month time.Month) time.Time {
	t.Helper()
	d := 1
	var hh, mm, ss, ns int
	return time.Date(year, month, d, hh, mm, ss, ns, time.Local)
}
