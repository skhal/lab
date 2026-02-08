// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csv_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/csv"
	"github.com/skhal/lab/x/fin/internal/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *pb.Date
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "january",
			s:    "2013.01",
			want: newDate(t, "Jan 2013"),
		},
		{
			name: "february",
			s:    "2013.02",
			want: newDate(t, "Feb 2013"),
		},
		{
			name: "march",
			s:    "2013.03",
			want: newDate(t, "Mar 2013"),
		},
		{
			name: "april",
			s:    "2013.04",
			want: newDate(t, "Apr 2013"),
		},
		{
			name: "may",
			s:    "2013.05",
			want: newDate(t, "May 2013"),
		},
		{
			name: "june",
			s:    "2013.06",
			want: newDate(t, "Jun 2013"),
		},
		{
			name: "july",
			s:    "2013.07",
			want: newDate(t, "Jul 2013"),
		},
		{
			name: "august",
			s:    "2013.08",
			want: newDate(t, "Aug 2013"),
		},
		{
			name: "september",
			s:    "2013.09",
			want: newDate(t, "Sep 2013"),
		},
		{
			name: "october",
			s:    "2013.1",
			want: newDate(t, "Oct 2013"),
		},
		{
			name: "november",
			s:    "2013.11",
			want: newDate(t, "Nov 2013"),
		},
		{
			name: "december",
			s:    "2013.12",
			want: newDate(t, "Dec 2013"),
		},
		{
			name:    "invalid month",
			s:       "2013.13",
			wantErr: true,
		},
		{
			name:    "negative month",
			s:       "2013.-1",
			wantErr: true,
		},
		{
			name:    "long month MMM",
			s:       "2013.123",
			wantErr: true,
		},
		{
			name:    "short year YYY",
			s:       "123.01",
			wantErr: true,
		},
		{
			name:    "long year YYYYY",
			s:       "12345.01",
			wantErr: true,
		},
		{
			name:    "negative year",
			s:       "-1.01",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := csv.ParseDate(tc.s)

			if err != nil {
				if !tc.wantErr {
					t.Fatalf("ParseDate(%q) failed: %v", tc.s, err)
				}
			} else if tc.wantErr {
				t.Fatalf("ParseDate(%q) succeeded unexpectedly", tc.s)
			}
			if diff := cmp.Diff(tc.want, d, protocmp.Transform()); diff != "" {
				t.Errorf("ParseDate(%q) mismatch (-want, +got):\n%s", tc.s, diff)
			}
		})
	}
}

func TestParseCents(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *pb.Cents
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "valid",
			s:    "123.45",
			want: newCents(t, 123, 45),
		},
		{
			name: "valid zero int part",
			s:    "0.45",
			want: newCents(t, 0, 45),
		},
		{
			name: "valid zero fractional part",
			s:    "123.00",
			want: newCents(t, 123, 0),
		},
		{
			name:    "must have integer part",
			s:       ".12",
			wantErr: true,
		},
		{
			name:    "short fractional part",
			s:       "1.1",
			wantErr: true,
		},
		{
			name:    "long fractional part",
			s:       "1.123",
			wantErr: true,
		},
		{
			name:    "no negative integer part",
			s:       "-1.12",
			wantErr: true,
		},
		{
			name:    "no negative fractional part",
			s:       "1.-2",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := csv.ParseCents(tc.s)

			if err != nil {
				if !tc.wantErr {
					t.Fatalf("ParseCents(%q) failed: %v", tc.s, err)
				}
			} else if tc.wantErr {
				t.Fatalf("ParseCents(%q) succeeded unexpectedly", tc.s)
			}
			if diff := cmp.Diff(tc.want, c, protocmp.Transform()); diff != "" {
				t.Errorf("ParseCents(%q) mismatch (-want, +got):\n%s", tc.s, diff)
			}
		})
	}
}

func newDate(t *testing.T, s string) *pb.Date {
	t.Helper()
	tm, err := time.Parse("Jan 2006", s)
	if err != nil {
		t.Fatal(err)
	}
	d := new(pb.Date)
	d.SetYear(int32(tm.Year()))
	d.SetMonth(int32(tm.Month()))
	return d
}

func newCents(t *testing.T, integerPart, fractionalPart int) *pb.Cents {
	t.Helper()
	c := new(pb.Cents)
	c.SetCents(int32(integerPart*100 + fractionalPart))
	return c
}
