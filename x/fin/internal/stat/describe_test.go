// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"testing"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/stat"
)

func TestDescribe_max(t *testing.T) {
	tests := []struct {
		name string
		cc   []fin.Cents
		want fin.Cents
	}{
		{
			name: "no values",
		},
		{
			name: "one value",
			cc:   []fin.Cents{1},
			want: 1,
		},
		{
			name: "two values ascending",
			cc:   []fin.Cents{1, 3},
			want: 3,
		},
		{
			name: "two values descending",
			cc:   []fin.Cents{3, 1},
			want: 3,
		},
		{
			name: "multiple values",
			cc:   []fin.Cents{2, 1, 3},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			d := stat.Describe(tc.cc)

			if d.Max != tc.want {
				t.Errorf("Describe(%s).Max = %s, want %s", tc.cc, d.Max, tc.want)
			}
		})
	}
}

func TestDescribe_min(t *testing.T) {
	tests := []struct {
		name string
		cc   []fin.Cents
		want fin.Cents
	}{
		{
			name: "no values",
		},
		{
			name: "one value",
			cc:   []fin.Cents{1},
			want: 1,
		},
		{
			name: "two values ascending",
			cc:   []fin.Cents{1, 3},
			want: 1,
		},
		{
			name: "two values descending",
			cc:   []fin.Cents{3, 1},
			want: 1,
		},
		{
			name: "multiple values",
			cc:   []fin.Cents{2, 1, 3},
			want: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			d := stat.Describe(tc.cc)

			if d.Min != tc.want {
				t.Errorf("Describe(%s).Min = %s, want %s", tc.cc, d.Min, tc.want)
			}
		})
	}
}

func TestDescribe_average(t *testing.T) {
	tests := []struct {
		name string
		cc   []fin.Cents
		want fin.Cents
	}{
		{
			name: "no values",
		},
		{
			name: "one value",
			cc:   []fin.Cents{1},
			want: 1,
		},
		{
			name: "two values ascending",
			cc:   []fin.Cents{1, 3},
			want: 2,
		},
		{
			name: "two values descending",
			cc:   []fin.Cents{3, 1},
			want: 2,
		},
		{
			name: "multiple values",
			cc:   []fin.Cents{2, 1, 3},
			want: 2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			d := stat.Describe(tc.cc)

			if d.Avg != tc.want {
				t.Errorf("Describe(%s).Avg = %s, want %s", tc.cc, d.Avg, tc.want)
			}
		})
	}
}

func TestDescribe_median(t *testing.T) {
	tests := []struct {
		name string
		cc   []fin.Cents
		want fin.Cents
	}{
		{
			name: "no values",
		},
		{
			name: "one value",
			cc:   []fin.Cents{1},
			want: 1,
		},
		{
			name: "two values ascending",
			cc:   []fin.Cents{1, 3},
			want: 2,
		},
		{
			name: "multiple values",
			cc:   []fin.Cents{1, 3, 7},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			d := stat.Describe(tc.cc)

			if d.Med != tc.want {
				t.Errorf("Describe(%s).Med = %s, want %s", tc.cc, d.Med, tc.want)
			}
		})
	}
}

func TestDescribe_stddev(t *testing.T) {
	tests := []struct {
		name string
		cc   []fin.Cents
		want fin.Cents
	}{
		{
			name: "no values",
		},
		{
			name: "one value",
			cc:   []fin.Cents{1},
		},
		{
			name: "two values ascending",
			cc:   []fin.Cents{1, 3},
			want: 1,
		},
		{
			name: "multiple values",
			cc:   []fin.Cents{1, 3, 7},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			d := stat.Describe(tc.cc)

			if d.Std != tc.want {
				t.Errorf("Describe(%s).Std = %s, want %s", tc.cc, d.Std, tc.want)
			}
		})
	}
}
