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
		},
		{
			name: "move only no points",
			p: func() plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: 11, Y: 12}},
				}
				return plot.Path{Commands: cmds}
			}(),
			want: "M 11 12",
		},
		{
			name: "move and points",
			p: func() plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: 11, Y: 12}},
					plot.PathLineCommand{Point: plot.Point{X: 21, Y: 22}},
					plot.PathLineCommand{Point: plot.Point{X: 31, Y: 32}},
				}
				return plot.Path{Commands: cmds}
			}(),
			want: "M 11 12 L 21 22 L 31 32",
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
