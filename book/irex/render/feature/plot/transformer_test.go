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

// floatRelativeDifferencePcent is a relative difference of two floating number.
const floatRelativeDifferencePcent = 0.001

func TestTransformer(t *testing.T) {
	tests := []struct {
		name  string
		x, y  float64
		tr    plot.Transformer
		wantX float64
		wantY float64
	}{
		{
			name:  "translate",
			x:     1,
			y:     2,
			tr:    plot.Translate(3, 4),
			wantX: 1 + 3,
			wantY: 2 + 4,
		},
		{
			name:  "scale",
			x:     1,
			y:     2,
			tr:    plot.Scale(3, 4),
			wantX: 1 * 3,
			wantY: 2 * 4,
		},
		{
			name:  "queue",
			x:     1,
			y:     2,
			tr:    plot.WithTransformer(plot.Translate(3, 4), plot.Scale(5, 6)),
			wantX: (1 + 3) * 5,
			wantY: (2 + 4) * 6,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x, y := tc.tr.Transform(tc.x, tc.y)

			opts := []cmp.Option{
				cmpopts.EquateApprox(floatRelativeDifferencePcent, 0),
			}
			if dx := cmp.Diff(tc.wantX, x, opts...); dx != "" {
				t.Errorf("X mismatch (-want +got):\n%s", dx)
			}
			if dy := cmp.Diff(tc.wantY, y, opts...); dy != "" {
				t.Errorf("Y mismatch (-want +got):\n%s", dy)
			}
		})
	}
}
