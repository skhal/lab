// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

func TestYearlyWithdrawer_Rebalance_illegalRate(t *testing.T) {
	tests := []struct {
		name    string
		rate    float64
		pos     fin.Position
		wantErr error
	}{
		{
			name:    "negative rate",
			rate:    -0.25,
			pos:     fin.Position{Investment: 123, Dividend: 56},
			wantErr: strategy.ErrWithdrawRate,
		},
		{
			name:    "above hundred percent",
			rate:    1.25,
			pos:     fin.Position{Investment: 123, Dividend: 56},
			wantErr: strategy.ErrWithdrawRate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			yw := strategy.YearlyWithdrawer{Rate: tc.rate}

			err := func(pos fin.Position) (err error) {
				defer func() {
					if x := recover(); x != nil {
						var ok bool
						if err, ok = x.(error); !ok {
							err = fmt.Errorf("test: recover: not an error: %v", x)
						}
					}
				}()
				for m := time.January; m <= time.December; m += 1 {
					pos = yw.Rebalance(pos)
				}
				pos = yw.Rebalance(pos)
				return nil
			}(tc.pos)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Rebalance() unexpected error %v, want %v", err, tc.wantErr)
				t.Log("Rate:", tc.rate)
			}
		})
	}
}

func TestYearlyWithdrawer_Rebalance_oneYear(t *testing.T) {
	tests := []struct {
		name string
		rate float64
		pos  fin.Position
		want fin.Position
	}{
		{
			name: "zero rate zero position",
		},
		{
			name: "zero rate",
			pos:  fin.Position{Investment: 123, Dividend: 56},
			want: fin.Position{Investment: 123, Dividend: 56},
		},
		{
			name: "zero rate no dividend",
			pos:  fin.Position{Investment: 123},
			want: fin.Position{Investment: 123},
		},
		{
			name: "zero rate no investment",
			pos:  fin.Position{Dividend: 123},
			want: fin.Position{Dividend: 123},
		},
		{
			name: "nonzero rate zero position",
			rate: 0.25,
		},
		{
			name: "nonzero rate",
			rate: 0.25,
			pos:  fin.Position{Investment: 123, Dividend: 56},
			want: deduct(t, fin.Position{Investment: 123, Dividend: 56}, 0.25),
		},
		{
			name: "nonzero rate no dividend",
			rate: 0.25,
			pos:  fin.Position{Investment: 123},
			want: deduct(t, fin.Position{Investment: 123}, 0.25),
		},
		{
			name: "nonzero rate no investment",
			rate: 0.25,
			pos:  fin.Position{Dividend: 123},
			want: deduct(t, fin.Position{Dividend: 123}, 0.25),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			yw := strategy.YearlyWithdrawer{Rate: tc.rate}

			// Must rebalance once every 12 months
			pos := tc.pos
			for m := time.January; m < time.December; m += 1 {
				pos = yw.Rebalance(pos)
			}
			got := yw.Rebalance(pos)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(*YearlyWithdrawer).Rebalance() mismatch (-want,+got):\n%s", diff)
				t.Log("Rate:", tc.rate)
				t.Log("Position:", tc.pos)
			}
		})
	}
}

func TestYearlyWithdrawer_Rebalance_twoYears(t *testing.T) {
	tests := []struct {
		name string
		rate float64
		pos  fin.Position
		want fin.Position
	}{
		{
			name: "zero rate zero position",
		},
		{
			name: "zero rate",
			pos:  fin.Position{Investment: 123, Dividend: 56},
			want: fin.Position{Investment: 123, Dividend: 56},
		},
		{
			name: "zero rate no dividend",
			pos:  fin.Position{Investment: 123},
			want: fin.Position{Investment: 123},
		},
		{
			name: "zero rate no investment",
			pos:  fin.Position{Dividend: 123},
			want: fin.Position{Dividend: 123},
		},
		{
			name: "nonzero rate zero position",
			rate: 0.25,
		},
		{
			name: "nonzero rate",
			rate: 0.25,
			pos:  fin.Position{Investment: 123, Dividend: 56},
			want: deduct(t, deduct(t, fin.Position{Investment: 123, Dividend: 56}, 0.25), 0.25),
		},
		{
			name: "nonzero rate no dividend",
			rate: 0.25,
			pos:  fin.Position{Investment: 123},
			want: deduct(t, deduct(t, fin.Position{Investment: 123}, 0.25), 0.25),
		},
		{
			name: "nonzero rate no investment",
			rate: 0.25,
			pos:  fin.Position{Dividend: 123},
			want: deduct(t, deduct(t, fin.Position{Dividend: 123}, 0.25), 0.25),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			yw := strategy.YearlyWithdrawer{Rate: tc.rate}

			// Must rebalance once every 12 months
			pos := tc.pos
			for y := 2; y > 0; y -= 1 { // 2 years
				for m := time.January; m <= time.December; m += 1 { // 12 months
					pos = yw.Rebalance(pos)
				}
			}
			got := yw.Rebalance(pos)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(*YearlyWithdrawer).Rebalance() mismatch (-want,+got):\n%s", diff)
				t.Log("Rate:", tc.rate)
				t.Log("Position:", tc.pos)
			}
		})
	}
}

func deduct(t *testing.T, pos fin.Position, rate float64) fin.Position {
	t.Helper()
	deduct := func(c fin.Cents) fin.Cents {
		return c - fin.Cents(math.Floor(float64(c)*rate))
	}
	pos.Investment = deduct(pos.Investment)
	pos.Dividend = deduct(pos.Dividend)
	return pos
}
