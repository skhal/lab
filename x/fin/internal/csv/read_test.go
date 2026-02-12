// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csv_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/csv"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/tests"
	"google.golang.org/protobuf/testing/protocmp"
)

var headers = []string{
	"h1date,h1sp,h1div,h1earn",
	"h2date,h2sp,h2div,h2earn",
	"h3date,h3sp,h3div,h3earn",
	"h4date,h4sp,h4div,h4earn",
	"h5date,h5sp,h5div,h5earn",
	"h6date,h6sp,h6div,h6earn",
	"h7date,h7sp,h7div,h7earn",
	"h8date,h8sp,h8div,h8earn",
}

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *pb.Market
		wantErr bool
	}{
		{
			name: "empty",
			want: newMarket(t),
		},
		{
			name: "header only",
			data: newCSV(t),
			want: newMarket(t),
		},
		{
			name: "one record",
			data: newCSV(t, "2006.01,1.11,1.22,1.33"),
			want: newMarket(t,
				tests.NewRecord(t, 2006, time.January, 111, 122, 133),
			),
		},
		{
			name: "two records",
			data: newCSV(t,
				"2006.01,1.11,1.22,1.33",
				"2006.02,2.11,2.22,2.33",
			),
			want: newMarket(t,
				tests.NewRecord(t, 2006, time.January, 111, 122, 133),
				tests.NewRecord(t, 2006, time.February, 211, 222, 233),
			),
		},
		{
			name: "two records end of year",
			data: newCSV(t,
				"2006.12,1.11,1.22,1.33",
				"2007.01,2.11,2.22,2.33",
			),
			want: newMarket(t,
				tests.NewRecord(t, 2006, time.December, 111, 122, 133),
				tests.NewRecord(t, 2007, time.January, 211, 222, 233),
			),
		},
		{
			name: "two records month gap",
			data: newCSV(t,
				"2006.01,1.11,1.22,1.33",
				"2006.03,2.11,2.22,2.33",
			),
			wantErr: true,
		},
		{
			name: "two records year gap",
			data: newCSV(t,
				"2006.01,1.11,1.22,1.33",
				"2007.02,2.11,2.22,2.33",
			),
			wantErr: true,
		},
		{
			name: "two records same month",
			data: newCSV(t,
				"2006.01,1.11,1.22,1.33",
				"2006.01,2.11,2.22,2.33",
			),
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.data)

			market, err := csv.Read(r)

			if tc.wantErr {
				if err == nil {
					t.Errorf("csv.Read() want error")
					t.Logf("data:\n%q", tc.data)
				}
			} else if err != nil {
				t.Errorf("csv.Read() unexpected error: %v", err)
				t.Logf("data:\n%q", tc.data)
			}
			if diff := cmp.Diff(tc.want, market, protocmp.Transform()); diff != "" {
				t.Errorf("csv.Read() mismatch (-want, +got):\n%s", diff)
				t.Logf("data:\n%q", tc.data)
			}
		})
	}
}

func newCSV(t *testing.T, data ...string) string {
	t.Helper()
	rows := slices.Concat(headers, data)
	return strings.Join(rows, "\n")
}

func newMarket(t *testing.T, recs ...*pb.Record) *pb.Market {
	t.Helper()
	return pb.Market_builder{
		Records: recs,
	}.Build()
}
