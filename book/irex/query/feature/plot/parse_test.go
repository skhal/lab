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

type testCase struct {
	name          string
	q             string
	wantSymbol    *pb.Symbol
	wantSinceDate *pb.Date
	wantUntilDate *pb.Date
	wantErr       error
}

func TestParse(t *testing.T) {
	tests := []testCase{
		{
			name:    "empty params",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name:    "invalid symbol",
			q:       "test-symbol",
			wantErr: plot.ErrNoSymbol,
		},
	}
	testParse(t, tests)
}

func TestParse_index(t *testing.T) {
	tests := []testCase{
		{
			name:       "spx",
			q:          "spx",
			wantSymbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
		},
		{
			name:    "multiple index",
			q:       "spx spx",
			wantErr: plot.ErrMultipleSymbol,
		},
	}
	testParse(t, tests)
}

func TestParse_indexMetric(t *testing.T) {
	tests := []testCase{
		{
			name:    "no index",
			q:       "div",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name:    "not index",
			q:       "cpi div",
			wantErr: plot.ErrNotIndex,
		},
		{
			name:    "multiple metrics",
			q:       "spx div earn",
			wantErr: plot.ErrMultipleIndexMetric,
		},
		{
			name:       "spx dividend",
			q:          "spx div",
			wantSymbol: newIndexMetricSymbol(t, pb.Symbol_Index_ID_SPX, pb.Symbol_Index_MET_DIV),
		},
		{
			name:       "spx earnings",
			q:          "spx earn",
			wantSymbol: newIndexMetricSymbol(t, pb.Symbol_Index_ID_SPX, pb.Symbol_Index_MET_EARN),
		},
	}
	testParse(t, tests)
}

func TestParse_marketMetric(t *testing.T) {
	tests := []testCase{
		{
			name:       "cpi",
			q:          "cpi",
			wantSymbol: newMarketMetricSymbol(t, pb.Symbol_Market_MET_CPI),
		},
		{
			name:    "multiple market metrics",
			q:       "cpi cpi",
			wantErr: plot.ErrMultipleSymbol,
		},
	}
	testParse(t, tests)
}

func TestParse_indexDate(t *testing.T) {
	tests := []testCase{
		{
			name:       "since misses date",
			q:          "spx since",
			wantSymbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
		},
		{
			name:    "since has no date",
			q:       "spx since abc",
			wantErr: plot.ErrSinceDate,
		},
		{
			name:          "with since date",
			q:             "spx since 1990-01",
			wantSymbol:    newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
			wantSinceDate: newDate(t, 1990, time.January, 31),
		},
		{
			name:       "until misses date",
			q:          "spx until",
			wantSymbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
		},
		{
			name:    "until has no date",
			q:       "spx until abc",
			wantErr: plot.ErrUntilDate,
		},
		{
			name:          "with until date",
			q:             "spx until 1990-01",
			wantSymbol:    newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
			wantUntilDate: newDate(t, 1990, time.January, 31),
		},
		{
			name:          "index with since and until dates",
			q:             "spx since 1990-01 until 1995-02",
			wantSymbol:    newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
			wantSinceDate: newDate(t, 1990, time.January, 31),
			wantUntilDate: newDate(t, 1995, time.February, 28),
		},
		{
			name:          "index with since and until dates in mixed order",
			q:             "until 1995-02 spx since 1990-01",
			wantSymbol:    newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
			wantSinceDate: newDate(t, 1990, time.January, 31),
			wantUntilDate: newDate(t, 1995, time.February, 28),
		},
	}
	testParse(t, tests)
}

func testParse(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := plot.Parse(tc.q)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%q unexpected error '%v'; want '%v'", tc.q, err, tc.wantErr)
			}
			var want *pb.PlotIntent
			if tc.wantSymbol != nil {
				want = newPlotIntent(t, tc.wantSymbol, tc.wantSinceDate, tc.wantUntilDate)
			}
			if d := cmp.Diff(want, got, protocmp.Transform()); d != "" {
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

func newIndexMetricSymbol(t *testing.T, idx pb.Symbol_Index_ID, m pb.Symbol_Index_Metric) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{
		Index: pb.Symbol_Index_builder{
			Id:     &idx,
			Metric: &m,
		}.Build(),
	}.Build()
}

func newMarketMetricSymbol(t *testing.T, m pb.Symbol_Market_Metric) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{
		Market: pb.Symbol_Market_builder{
			Metric: &m,
		}.Build(),
	}.Build()
}
