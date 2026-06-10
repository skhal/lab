// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/render/feature/plot"
)

func TestPlotter_Plot(t *testing.T) {
	const (
		xrange = 200
		yrange = 100
	)
	type plotter interface {
		Plot([]*pb.PlotFeature_Quote) *plot.Path
	}
	tests := []struct {
		name   string
		pl     plotter
		quotes []*pb.PlotFeature_Quote
		want   *plot.Path
	}{
		{
			name: "no quotes",
			pl:   plot.NewPlotter(xrange, yrange),
		},
		{
			name: "one quote",
			pl:   plot.NewPlotter(xrange, yrange),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
			},
			// place in the middle of the plot
			want: func() *plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: xrange / 2, Y: yrange / 2}},
				}
				return &plot.Path{Commands: cmds}
			}(),
		},
		{
			name: "two quotes ascend",
			pl:   plot.NewPlotter(xrange, yrange),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 102),
			},
			// place in the opposite corners of the plot
			want: func() *plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: 0, Y: 0}},
					plot.PathLineCommand{Point: plot.Point{X: 200, Y: 100}},
				}
				return &plot.Path{Commands: cmds}
			}(),
		},
		{
			name: "two quotes descend",
			pl:   plot.NewPlotter(xrange, yrange),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 100),
			},
			// place in the opposite corners of the plot
			want: func() *plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: 0, Y: 100}},
					plot.PathLineCommand{Point: plot.Point{X: 200, Y: 0}},
				}
				return &plot.Path{Commands: cmds}
			}(),
		},
		{
			name: "three quotes",
			pl:   plot.NewPlotter(xrange, yrange),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 102),
				newQuote(t, 1990, time.March, 31, 103),
			},
			want: func() *plot.Path {
				cmds := []plot.PathCommand{
					plot.PathMoveCommand{Point: plot.Point{X: 0, Y: 0}},
					plot.PathLineCommand{Point: plot.Point{X: 100, Y: 50}},
					plot.PathLineCommand{Point: plot.Point{X: 200, Y: 100}},
				}
				return &plot.Path{Commands: cmds}
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.pl.Plot(tc.quotes)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(plot.PathMoveCommand{}, plot.PathLineCommand{}),
			}
			if d := cmp.Diff(tc.want, got, opts...); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Log(tc.quotes)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, cent int32) *pb.PlotFeature_Quote {
	t.Helper()
	return pb.PlotFeature_Quote_builder{
		Date: newDate(t, year, month, day),
		Cent: pb.Cent_builder{Value: &cent}.Build(),
	}.Build()
}

func newDate(t *testing.T, year int32, month time.Month, day int32) *pb.Date {
	t.Helper()
	return pb.Date_builder{
		Year:  &year,
		Month: new(int32(month)),
		Day:   &day,
	}.Build()
}
