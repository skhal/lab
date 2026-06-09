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
		name string
		x    float64
		tr   *plot.Transformer
		want float64
	}{
		{
			name: "zero offset non-zero scale",
			x:    1,
			tr:   plot.NewTransformer(0, 3),
			want: func() float64 { return (1 - 0) * 3 }(),
		},
		{
			name: "non-zero offset zero scale",
			x:    1,
			tr:   plot.NewTransformer(2, 0),
			want: func() float64 { return (1 - 2) * 0 }(),
		},
		{
			name: "non-zero offset non-zero scale",
			x:    1,
			tr:   plot.NewTransformer(2, 3),
			want: func() float64 { return (1 - 2) * 3 }(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.tr.Transform(tc.x)

			opts := []cmp.Option{
				cmpopts.EquateApprox(floatRelativeDifferencePcent, 0),
			}
			if d := cmp.Diff(tc.want, got, opts...); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Logf("%#v", tc.tr)
			}
		})
	}
}
