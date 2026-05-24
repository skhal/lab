// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport_test

import (
	"errors"
	"fmt"
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
			rec:  []string{"1990.01", "1.01", "1.02"},
			want: newQuote(t, 1990, time.January, 31, 101, 102),
		},
	}
	testParse(t, tests)
}

func TestParse_date(t *testing.T) {
	tests := []parseTest{
		{
			name:    "empty",
			rec:     []string{"", "1.01", "1.02"},
			wantErr: csvimport.ErrDate,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.a1", "1.01", "1.02"},
			wantErr: csvimport.ErrDate,
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
			rec:     []string{"1990.01", "", "1.02"},
			wantErr: csvimport.ErrSPX,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.01", "1.a1", "1.02"},
			wantErr: csvimport.ErrSPX,
		},
	}
	testParse(t, tests)
}

func TestParse_div(t *testing.T) {
	tests := []parseTest{
		{
			name:    "no field",
			rec:     []string{"1990.01", "1.01"},
			wantErr: csvimport.ErrNoDividend,
		},
		{
			name:    "empty field",
			rec:     []string{"1990.01", "1.01", ""},
			wantErr: csvimport.ErrDividend,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.01", "1.01", "1.a2"},
			wantErr: csvimport.ErrDividend,
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

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    *pb.Date
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: csvimport.ErrFormat,
		},
		{
			name:    "no year",
			value:   ".01",
			wantErr: csvimport.ErrYear,
		},
		{
			name:    "invalid year",
			value:   fmt.Sprintf("%d.01", csvimport.MinYear-1),
			wantErr: csvimport.ErrYear,
		},
		{
			name:    "no month",
			value:   "1990.",
			wantErr: csvimport.ErrMonth,
		},
		{
			name:    "invalid month",
			value:   "1990.13",
			wantErr: csvimport.ErrMonth,
		},
		{
			// single digit is only for October (10) that is .10 become .1
			name:    "single digit month only october",
			value:   "1990.2",
			wantErr: csvimport.ErrMonth,
		},
		{
			name:  "january",
			value: "1990.01",
			want:  newDate(t, 1990, time.January, 31),
		},
		{
			name:  "february",
			value: "1990.02",
			want:  newDate(t, 1990, time.February, 28),
		},
		{
			name:  "march",
			value: "1990.03",
			want:  newDate(t, 1990, time.March, 31),
		},
		{
			name:  "april",
			value: "1990.04",
			want:  newDate(t, 1990, time.April, 30),
		},
		{
			name:  "may",
			value: "1990.05",
			want:  newDate(t, 1990, time.May, 31),
		},
		{
			name:  "june",
			value: "1990.06",
			want:  newDate(t, 1990, time.June, 30),
		},
		{
			name:  "july",
			value: "1990.07",
			want:  newDate(t, 1990, time.July, 31),
		},
		{
			name:  "august",
			value: "1990.08",
			want:  newDate(t, 1990, time.August, 31),
		},
		{
			name:  "september",
			value: "1990.09",
			want:  newDate(t, 1990, time.September, 30),
		},
		{
			name:  "october",
			value: "1990.1",
			want:  newDate(t, 1990, time.October, 31),
		},
		{
			name:  "november",
			value: "1990.11",
			want:  newDate(t, 1990, time.November, 30),
		},
		{
			name:  "december",
			value: "1990.12",
			want:  newDate(t, 1990, time.December, 31),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := csvimport.ParseDate(tc.value)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("ParseCent(%q) mismatch (-want +got):\n%s", tc.value, d)
			}
		})
	}
}

func TestParseCent(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    *pb.Cent
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: csvimport.ErrFormat,
		},
		{
			name:    "no dollar",
			value:   ".02",
			wantErr: csvimport.ErrDollar,
		},
		{
			name:    "dollar not number",
			value:   "a.23",
			wantErr: csvimport.ErrDollar,
		},
		{
			name:    "no cent",
			value:   "1.",
			wantErr: csvimport.ErrCent,
		},
		{
			name:    "one digit cent",
			value:   "1.2",
			wantErr: csvimport.ErrCent,
		},
		{
			name:  "two digit cent",
			value: "1.23",
			want:  pb.Cent_builder{Value: new(int32(123))}.Build(),
		},
		{
			name:    "three digit cent",
			value:   "1.234",
			wantErr: csvimport.ErrCent,
		},
		{
			name:    "cent not number",
			value:   "1.ab",
			wantErr: csvimport.ErrCent,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := csvimport.ParseCent(tc.value)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("ParseCent(%q) mismatch (-want +got):\n%s", tc.value, d)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx, div int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: newDate(t, year, month, day),
		Spx:  pb.Cent_builder{Value: &spx}.Build(),
		Div:  pb.Cent_builder{Value: &div}.Build(),
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
