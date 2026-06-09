// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"

	"github.com/skhal/lab/book/irex/render/feature/plot"
)

func TestViewBox_String(t *testing.T) {
	tests := []struct {
		name string
		vb   plot.ViewBox
		want string
	}{
		{
			name: "empty",
			want: "0 0 0 0",
		},
		{
			name: "zero origin",
			vb:   plot.ViewBox{Width: 3, Height: 4},
			want: "0 0 3 4",
		},
		{
			name: "zero dimension",
			vb:   plot.ViewBox{Point: plot.Point{1, 2}},
			want: "1 2 0 0",
		},
		{
			name: "non zero",
			vb:   plot.ViewBox{Point: plot.Point{1, 2}, Width: 3, Height: 4},
			want: "1 2 3 4",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.vb.String()

			if got != tc.want {
				t.Errorf("String() = %q; want %q", got, tc.want)
			}
		})
	}
}
