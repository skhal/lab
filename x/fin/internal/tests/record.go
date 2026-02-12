// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"testing"
	"time"

	"github.com/skhal/lab/x/fin/internal/pb"
)

// NewRecord creates a market record for tests.
func NewRecord(t *testing.T, y int32, m time.Month, sp, div, earn int32) *pb.Record {
	t.Helper()
	month := int32(m)
	return pb.Record_builder{
		Date: pb.Date_builder{
			Year:  &y,
			Month: &month,
		}.Build(),
		Quote: pb.Quote_builder{
			SpComposite: pb.Cents_builder{Cents: &sp}.Build(),
			Dividend:    pb.Cents_builder{Cents: &div}.Build(),
			Earnings:    pb.Cents_builder{Cents: &earn}.Build(),
		}.Build(),
	}.Build()
}
