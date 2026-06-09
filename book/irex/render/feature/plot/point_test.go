// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"

	"github.com/skhal/lab/book/irex/render/feature/plot"
)

func TestPoint(t *testing.T) {
	tests := []struct {
		name string
		p    plot.Point
		want string
	}{
		{
			name: "empty",
			want: "0 0",
		},
		{
			name: "x only",
			p:    plot.Point{X: 1},
			want: "1 0",
		},
		{
			name: "y only",
			p:    plot.Point{Y: 1},
			want: "0 1",
		},
		{
			name: "x and y",
			p:    plot.Point{X: 1, Y: 2},
			want: "1 2",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.p.String()

			if s != tc.want {
				t.Errorf("Point{}.String() = %q; want %q", s, tc.want)
				t.Logf("%#v", tc.p)
			}
		})
	}
}
