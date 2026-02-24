// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report_test

import (
	"strings"
	"testing"

	gotests "github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/fin/internal/report"
	"github.com/skhal/lab/x/fin/internal/sim"
	"github.com/skhal/lab/x/fin/internal/stat"
)

func TestSimulation(t *testing.T) {
	tests := []struct {
		name   string
		data   []*report.SimulationData
		golden gotests.GoldenFile
	}{
		{
			name:   "empty",
			golden: gotests.GoldenFile("testdata/simulation/empty.txt"),
		},
		{
			name: "one strategy",
			data: []*report.SimulationData{
				{
					Name: "test",
					Result: &sim.Result{
						Start: 10,
						End: stat.Description{
							Max: 11,
							Min: 12,
							Avg: 13,
							Med: 14,
							Std: 15,
						},
					},
				},
			},
			golden: gotests.GoldenFile("testdata/simulation/one_strategy.txt"),
		},
		{
			name: "two strategies",
			data: []*report.SimulationData{
				{
					Name: "test-one",
					Result: &sim.Result{
						Start: 10,
						End: stat.Description{
							Max: 11,
							Min: 12,
							Avg: 13,
							Med: 14,
							Std: 15,
						},
					},
				},
				{
					Name: "test-two",
					Result: &sim.Result{
						Start: 20,
						End: stat.Description{
							Max: 21,
							Min: 22,
							Avg: 23,
							Med: 24,
							Std: 25,
						},
					},
				},
			},
			golden: gotests.GoldenFile("testdata/simulation/two_strategies.txt"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder

			err := report.Simulation(&b, tc.data)

			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got := b.String()
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("Simulation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
