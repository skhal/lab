// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

type parseTest struct {
	name    string
	rec     []string
	want    *pb.Quote
	wantErr error
}

func TestParse(t *testing.T) {
	tests := []parseTest{
		{
			name:    "empty",
			wantErr: csvimport.ErrNoDate,
		},
		{
			name: "valid record",
			rec:  []string{"1990.01", "1.23"},
			want: newQuote(t, 1990, time.January, 31, 123),
		},
	}
	testParse(t, tests)
}

func TestParse_date(t *testing.T) {
	tests := []parseTest{
		{
			name:    "empty date",
			rec:     []string{"", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name:    "no year",
			rec:     []string{".01", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name:    "invalid year",
			rec:     []string{"1600", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name:    "no month",
			rec:     []string{"1990.", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name:    "invalid month",
			rec:     []string{"1990.13", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name: "single digit month only october",
			// single digit is only for October (10) that is .10 become .1
			rec:     []string{"1990.2", "1.23"},
			wantErr: csvimport.ErrDate,
		},
		{
			name: "january",
			rec:  []string{"1990.01", "1.23"},
			want: newQuote(t, 1990, time.January, 31, 123),
		},
		{
			name: "february",
			rec:  []string{"1990.02", "1.23"},
			want: newQuote(t, 1990, time.February, 28, 123),
		},
		{
			name: "march",
			rec:  []string{"1990.03", "1.23"},
			want: newQuote(t, 1990, time.March, 31, 123),
		},
		{
			name: "april",
			rec:  []string{"1990.04", "1.23"},
			want: newQuote(t, 1990, time.April, 30, 123),
		},
		{
			name: "may",
			rec:  []string{"1990.05", "1.23"},
			want: newQuote(t, 1990, time.May, 31, 123),
		},
		{
			name: "june",
			rec:  []string{"1990.06", "1.23"},
			want: newQuote(t, 1990, time.June, 30, 123),
		},
		{
			name: "july",
			rec:  []string{"1990.07", "1.23"},
			want: newQuote(t, 1990, time.July, 31, 123),
		},
		{
			name: "august",
			rec:  []string{"1990.08", "1.23"},
			want: newQuote(t, 1990, time.August, 31, 123),
		},
		{
			name: "september",
			rec:  []string{"1990.09", "1.23"},
			want: newQuote(t, 1990, time.September, 30, 123),
		},
		{
			name: "october",
			rec:  []string{"1990.1", "1.23"},
			want: newQuote(t, 1990, time.October, 31, 123),
		},
		{
			name: "november",
			rec:  []string{"1990.11", "1.23"},
			want: newQuote(t, 1990, time.November, 30, 123),
		},
		{
			name: "december",
			rec:  []string{"1990.12", "1.23"},
			want: newQuote(t, 1990, time.December, 31, 123),
		},
	}
	testParse(t, tests)
}

func TestParse_spx(t *testing.T) {
	tests := []parseTest{
		{
			name:    "no field",
			rec:     []string{"1990.01"},
			wantErr: csvimport.ErrNoSPX,
		},
		{
			name:    "empty field",
			rec:     []string{"1990.01", ""},
			wantErr: csvimport.ErrSPX,
		},
		{
			name:    "no dollar",
			rec:     []string{"1990.01", ".12"},
			wantErr: csvimport.ErrSPX,
		},
		{
			name:    "no cents",
			rec:     []string{"1990.01", "1."},
			wantErr: csvimport.ErrSPX,
		},
		{
			name:    "two digit cents",
			rec:     []string{"1990.01", "1.2"},
			wantErr: csvimport.ErrSPX,
		},
		{
			name:    "cents not number",
			rec:     []string{"1990.01", "1.ab"},
			wantErr: csvimport.ErrSPX,
		},
	}
	testParse(t, tests)
}

func testParse(t *testing.T, tests []parseTest) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := csvimport.Parse(tc.rec)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Logf("record:\n%v", tc.rec)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: pb.Date_builder{
			Year:  &year,
			Month: new(int32(month)),
			Day:   &day,
		}.Build(),
		Spx: pb.Cent_builder{Value: &spx}.Build(),
	}.Build()
}
