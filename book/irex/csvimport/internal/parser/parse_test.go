// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport/internal/parser"
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
			wantErr: parser.ErrNoDate,
		},
		{
			name: "valid record",
			rec:  []string{"1990.01", "1.01", "1.02", "1.03"},
			want: newQuote(t, 1990, time.January, 31, 101, 102, 103),
		},
	}
	testParse(t, tests)
}

func TestParse_date(t *testing.T) {
	tests := []parseTest{
		{
			name:    "empty",
			rec:     []string{"", "1.01", "1.02"},
			wantErr: parser.ErrDate,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.a1", "1.01", "1.02"},
			wantErr: parser.ErrDate,
		},
	}
	testParse(t, tests)
}

func TestParse_spx(t *testing.T) {
	tests := []parseTest{
		{
			name:    "no field",
			rec:     []string{"1990.01"},
			wantErr: parser.ErrNoSPX,
		},
		{
			name:    "empty field",
			rec:     []string{"1990.01", "", "1.02"},
			wantErr: parser.ErrSPX,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.01", "1.a1", "1.02"},
			wantErr: parser.ErrSPX,
		},
	}
	testParse(t, tests)
}

func TestParse_div(t *testing.T) {
	tests := []parseTest{
		{
			name:    "no field",
			rec:     []string{"1990.01", "1.01"},
			wantErr: parser.ErrNoDividend,
		},
		{
			name:    "empty field",
			rec:     []string{"1990.01", "1.01", ""},
			wantErr: parser.ErrDividend,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.01", "1.01", "1.a2"},
			wantErr: parser.ErrDividend,
		},
	}
	testParse(t, tests)
}

func TestParse_cpi(t *testing.T) {
	tests := []parseTest{
		{
			name:    "no field",
			rec:     []string{"1990.01", "1.01", "1.02"},
			wantErr: parser.ErrNoCPI,
		},
		{
			name:    "empty field",
			rec:     []string{"1990.01", "1.01", "1.02", ""},
			wantErr: parser.ErrCPI,
		},
		{
			name:    "invalid",
			rec:     []string{"1990.01", "1.01", "1.02", "1.a3"},
			wantErr: parser.ErrCPI,
		},
	}
	testParse(t, tests)
}

func testParse(t *testing.T, tests []parseTest) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.ParseQuote(tc.rec)

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
			wantErr: parser.ErrFormat,
		},
		{
			name:    "no year",
			value:   ".01",
			wantErr: parser.ErrYear,
		},
		{
			name:    "invalid year",
			value:   fmt.Sprintf("%d.01", parser.MinYear-1),
			wantErr: parser.ErrYear,
		},
		{
			name:    "no month",
			value:   "1990.",
			wantErr: parser.ErrMonth,
		},
		{
			name:    "invalid month",
			value:   "1990.13",
			wantErr: parser.ErrMonth,
		},
		{
			// single digit is only for October (10) that is .10 become .1
			name:    "single digit month only october",
			value:   "1990.2",
			wantErr: parser.ErrMonth,
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
			got, err := parser.ParseDate(tc.value)

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
			wantErr: parser.ErrFormat,
		},
		{
			name:    "no dollar",
			value:   ".02",
			wantErr: parser.ErrDollar,
		},
		{
			name:    "dollar not number",
			value:   "a.23",
			wantErr: parser.ErrDollar,
		},
		{
			name:    "no cent",
			value:   "1.",
			wantErr: parser.ErrCent,
		},
		{
			name:    "one digit cent",
			value:   "1.2",
			wantErr: parser.ErrCent,
		},
		{
			name:  "two digit cent",
			value: "1.23",
			want:  pb.Cent_builder{Value: new(int32(123))}.Build(),
		},
		{
			name:    "three digit cent",
			value:   "1.234",
			wantErr: parser.ErrCent,
		},
		{
			name:    "cent not number",
			value:   "1.ab",
			wantErr: parser.ErrCent,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.ParseCent(tc.value)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("ParseCent(%q) mismatch (-want +got):\n%s", tc.value, d)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx, div, cpi int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: newDate(t, year, month, day),
		Spx:  pb.Cent_builder{Value: &spx}.Build(),
		Div:  pb.Cent_builder{Value: &div}.Build(),
		Cpi:  pb.Cent_builder{Value: &cpi}.Build(),
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
