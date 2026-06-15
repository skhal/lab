// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query/feature/plot"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParse(t *testing.T) {
	var nilDate *pb.Date = nil
	tests := []struct {
		name    string
		q       string
		want    *pb.PlotIntent
		wantErr error
	}{
		{
			name:    "empty params",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name:    "invalid symbol",
			q:       "test-symbol",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name: "spx",
			q:    "spx",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), nilDate, nilDate),
		},
		{
			name: "since misses date",
			q:    "spx since",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), nilDate, nilDate),
		},
		{
			name:    "since has no date",
			q:       "spx since abc",
			wantErr: plot.ErrSinceDate,
		},
		{
			name: "with since date",
			q:    "spx since 1990-01",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), newDate(t, 1990, time.January, 31), nilDate),
		},
		{
			name: "until misses date",
			q:    "spx until",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), nilDate, nilDate),
		},
		{
			name:    "until has no date",
			q:       "spx until abc",
			wantErr: plot.ErrUntilDate,
		},
		{
			name: "with until date",
			q:    "spx until 1990-01",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), nilDate, newDate(t, 1990, time.January, 31)),
		},
		{
			name: "index with since and until dates",
			q:    "spx since 1990-01 until 1995-02",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), newDate(t, 1990, time.January, 31), newDate(t, 1995, time.February, 28)),
		},
		{
			name: "index with since and until dates in mixed order",
			q:    "until 1995-02 spx since 1990-01",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_Index_ID_SPX), newDate(t, 1990, time.January, 31), newDate(t, 1995, time.February, 28)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := plot.Parse(tc.q)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%q unexpected error '%v'; want '%v'", tc.q, err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("%q mismatch (-want +got):\n%s", tc.q, d)
			}
		})
	}
}

func newPlotIntent(t *testing.T, s *pb.Symbol, since, until *pb.Date) *pb.PlotIntent {
	t.Helper()
	return pb.PlotIntent_builder{
		Symbol: s,
		Since:  since,
		Until:  until,
	}.Build()
}

func newIndexSymbol(t *testing.T, idx pb.Symbol_Index_ID) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{
		Index: pb.Symbol_Index_builder{
			Id: &idx,
		}.Build(),
	}.Build()
}

func newDate(t *testing.T, year int32, month time.Month, day int32) *pb.Date {
	t.Helper()
	return pb.Date_builder{
		Year:  &year,
		Month: new(int32(month)),
		Day:   &day,
	}.Build()
}
