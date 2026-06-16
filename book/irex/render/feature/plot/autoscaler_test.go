// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/book/irex/render/feature/plot"
)

const floatRelDiffPcent = 0.001 // 0.1%

func TestAutoscaler_Scale(t *testing.T) {
	tests := []struct {
		name       string
		padding    plot.Percent
		xmin, xmax float64
		wantMin    float64
		wantMax    float64
	}{
		{
			name:    "no padding",
			xmin:    2,
			xmax:    4,
			wantMin: 2,
			wantMax: 4,
		},
		{
			name:    "scale max",
			padding: plot.Percent(0.01),
			xmax:    4,
			wantMin: 0,
			wantMax: func() float64 { return 4 + (4-0)*.01 }(),
		},
		{
			name:    "scale min and max",
			padding: plot.Percent(0.01),
			xmin:    2,
			xmax:    4,
			wantMin: func() float64 { return max(0, 2-(4-2)*.01) }(),
			wantMax: func() float64 { return 4 + (4-2)*.01 }(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scaler := plot.Autoscaler{tc.padding}

			xmin, xmax := scaler.Scale(tc.xmin, tc.xmax)

			opts := []cmp.Option{
				cmpopts.EquateApprox(floatRelDiffPcent, 0),
			}
			if d := cmp.Diff(tc.wantMin, xmin, opts...); d != "" {
				t.Errorf("xmin mismatch (-want +got):\n%s", d)
			}
			if d := cmp.Diff(tc.wantMax, xmax, opts...); d != "" {
				t.Errorf("xmax mismatch (-want +got):\n%s", d)
			}
		})
	}
}
