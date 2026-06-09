// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"

	"github.com/skhal/lab/book/irex/render/feature/plot"
)

func TestPath(t *testing.T) {
	tests := []struct {
		name string
		p    plot.Path
		want string
	}{
		{
			name: "empty",
			want: "M 0 0",
		},
		{
			name: "move only no points",
			p:    plot.Path{Move: plot.Point{X: 11, Y: 12}},
			want: "M 11 12",
		},
		{
			name: "move and points",
			p: plot.Path{Move: plot.Point{X: 11, Y: 12}, Line: []plot.Point{
				{21, 22}, {31, 32},
			}},
			want: "M 11 12 L 21 22 L 31 32",
		},
		{
			name: "points only no move",
			p: plot.Path{Line: []plot.Point{
				{21, 22}, {31, 32},
			}},
			want: "M 0 0 L 21 22 L 31 32",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.p.String()

			if s != tc.want {
				t.Errorf("Path{}.String() = %q; want %q", s, tc.want)
				t.Logf("%#v", tc.p)
			}
		})
	}
}
