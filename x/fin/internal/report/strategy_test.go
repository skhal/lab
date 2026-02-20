// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report_test

import (
	"flag"
	"fmt"
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

func TestStrategies(t *testing.T) {
	tt := []struct {
		name    string
		infos   []*report.StrategyInfo
		golden  gotests.GoldenFile
		wantErr bool
	}{
		{
			name: "one strategy",
			infos: []*report.StrategyInfo{
				{
					Name:        "test-strategy",
					Description: "test strategy description",
					Start:       tests.NewBalance(t, 2006, time.January, 100),
					End:         tests.NewBalance(t, 2006, time.February, 110),
				},
			},
			golden: gotests.GoldenFile("testdata/strategies_one.txt"),
		},
		{
			name: "two strategies",
			infos: []*report.StrategyInfo{
				{
					Name:        "test-strategy-a",
					Description: "test strategy A description",
					Start:       tests.NewBalance(t, 2006, time.January, 100),
					End:         tests.NewBalance(t, 2006, time.February, 110),
				},
				{
					Name:        "test-strategy-b",
					Description: "test strategy B description",
					Start:       tests.NewBalance(t, 2006, time.January, 100),
					End:         tests.NewBalance(t, 2006, time.February, 120),
				},
			},
			golden: gotests.GoldenFile("testdata/strategies_two.txt"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := new(strings.Builder)

			err := report.Strategies(b, tc.infos)
			got := b.String()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("report.Strategies() want error\n%s", infosStringer(tc.infos))
				}
			} else {
				if err != nil {
					t.Fatalf("report.Strategies() unexpected error %v\n%s", err, infosStringer(tc.infos))
				}
			}
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("report.Strategies() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

type infosStringer []*report.StrategyInfo

func (infos infosStringer) String() string {
	b := new(strings.Builder)
	for i, info := range infos {
		if i > 0 {
			fmt.Fprintln(b)
		}
		fmt.Fprint(b, info)
	}
	return b.String()
}
