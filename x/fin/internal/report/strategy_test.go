// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report_test

import (
	"flag"
	"strings"
	"testing"
	"time"

	gotests "github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/fin/internal/report"
	"github.com/skhal/lab/x/fin/internal/tests"
)

var update = flag.Bool("update", false, "update golden files")

func TestStrategy(t *testing.T) {
	tt := []struct {
		name    string
		info    report.StrategyInfo
		golden  gotests.GoldenFile
		wantErr bool
	}{
		{
			name: "strategy",
			info: report.StrategyInfo{
				Name:        "test-strategy",
				Description: "test strategy description",
				Start:       tests.NewBalance(t, 2006, time.January, 100),
				End:         tests.NewBalance(t, 2006, time.February, 110),
			},
			golden: gotests.GoldenFile("testdata/strategy.txt"),
		},
		{
			name: "no end balance",
			info: report.StrategyInfo{
				Name:        "test-strategy",
				Description: "test strategy description",
				Start:       tests.NewBalance(t, 2006, time.January, 100),
			},
			golden: gotests.GoldenFile("testdata/strategy_no_end_balance.txt"),
		},
		{
			name: "no start balance",
			info: report.StrategyInfo{
				Name:        "test-strategy",
				Description: "test strategy description",
				End:         tests.NewBalance(t, 2006, time.January, 100),
			},
			golden: gotests.GoldenFile("testdata/strategy_no_start_balance.txt"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := new(strings.Builder)

			err := report.Strategy(b, &tc.info)
			got := b.String()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("report.Strategy() want error\n%s", tests.InfoStringer(tc.info))
				}
			} else {
				if err != nil {
					t.Fatalf("report.Strategy() unexpected error %v\n%s", err, tests.InfoStringer(tc.info))
				}
			}
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("report.Strategy() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
