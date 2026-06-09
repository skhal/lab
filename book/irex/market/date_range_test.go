// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package market_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/market"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestDateRange_Range(t *testing.T) {
	tests := []struct {
		name   string
		quotes []*pb.Quote
		since  *pb.Date
		until  *pb.Date
		want   []*pb.Quote
	}{
		{
			name: "no date limits",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
		},
		{
			name: "since date",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			since: newDate(t, 1990, time.February, 28),
			want: []*pb.Quote{
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
		},
		{
			name: "until date",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			until: newDate(t, 1990, time.February, 28),
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
			},
		},
		{
			name: "range dates",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			since: newDate(t, 1990, time.February, 28),
			until: newDate(t, 1990, time.March, 31),
			want: []*pb.Quote{
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
			},
		},
		{
			name: "empty range",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			since: newDate(t, 1990, time.March, 31),
			until: newDate(t, 1990, time.February, 28),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dr := market.DateRange{Since: tc.since, Until: tc.until}

			got := slices.Collect(dr.Quotes(tc.quotes))

			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}
