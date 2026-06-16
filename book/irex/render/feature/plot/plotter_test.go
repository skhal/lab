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
		xrange           = 200
		yrange           = 100
		equateApproxFrac = plot.PlotPaddingPercent
	)
	type plotter interface {
		Plot([]*pb.PlotFeature_Quote) *plot.PlotInfo
	}
	tests := []struct {
		name   string
		pl     plotter
		quotes []*pb.PlotFeature_Quote
		want   *plot.PlotInfo
	}{
		{
			name: "no quotes",
			pl:   plot.NewPlotter(plot.Scale(xrange, yrange)),
		},
		{
			name: "one quote",
			pl:   plot.NewPlotter(plot.Scale(xrange, yrange)),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
			},
			want: &plot.PlotInfo{
				// place the quote in the middle of the plot
				Path: func() *plot.Path {
					cmds := []plot.PathCommand{
						plot.PathMoveCommand{
							Point: plot.Point{X: xrange / 2, Y: yrange / 2},
						},
					}
					return &plot.Path{Commands: cmds}
				}(),
				Quotes: plot.XQuote{
					xrange / 2: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.January, 31),
						Cents:    101,
					},
				},
			},
		},
		{
			name: "two quotes ascending",
			pl:   plot.NewPlotter(plot.Scale(xrange, yrange)),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 102),
			},
			want: &plot.PlotInfo{
				Path: func() *plot.Path {
					// place in the opposite corners of the plot
					cmds := []plot.PathCommand{
						plot.PathMoveCommand{Point: plot.Point{X: 0, Y: 0}},
						plot.PathLineCommand{Point: plot.Point{X: xrange, Y: yrange}},
					}
					return &plot.Path{Commands: cmds}
				}(),
				Quotes: plot.XQuote{
					0: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.January, 31),
						Cents:    101,
					},
					xrange: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.February, 28),
						Cents:    102,
					},
				},
			},
		},
		{
			name: "two quotes descending",
			pl:   plot.NewPlotter(plot.Scale(xrange, yrange)),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 100),
			},
			// place in the opposite corners of the plot
			want: &plot.PlotInfo{
				Path: func() *plot.Path {
					cmds := []plot.PathCommand{
						plot.PathMoveCommand{Point: plot.Point{X: 0, Y: yrange}},
						plot.PathLineCommand{Point: plot.Point{X: xrange, Y: 0}},
					}
					return &plot.Path{Commands: cmds}
				}(),
				Quotes: plot.XQuote{
					0: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.January, 31),
						Cents:    101,
					},
					xrange: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.February, 28),
						Cents:    100,
					},
				},
			},
		},
		{
			name: "three quotes",
			pl:   plot.NewPlotter(plot.Scale(xrange, yrange)),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 102),
				newQuote(t, 1990, time.March, 31, 103),
			},
			want: &plot.PlotInfo{
				Path: func() *plot.Path {
					cmds := []plot.PathCommand{
						plot.PathMoveCommand{Point: plot.Point{X: 0, Y: 0}},
						plot.PathLineCommand{Point: plot.Point{X: xrange / 2, Y: yrange / 2}},
						plot.PathLineCommand{Point: plot.Point{X: xrange, Y: yrange}},
					}
					return &plot.Path{Commands: cmds}
				}(),
				Quotes: plot.XQuote{
					0: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.January, 31),
						Cents:    101,
					},
					xrange / 2: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.February, 28),
						Cents:    102,
					},
					xrange: &plot.Quote{
						UnixTime: unixTime(t, 1990, time.March, 31),
						Cents:    103,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info := tc.pl.Plot(tc.quotes)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(plot.PathMoveCommand{}, plot.PathLineCommand{}),
				cmpopts.IgnoreFields(plot.PathMoveCommand{}, "Y"),
				cmpopts.IgnoreFields(plot.PathLineCommand{}, "Y"),
				cmpopts.IgnoreFields(plot.PlotInfo{}, "Ymin", "Ymax"),
				cmpopts.EquateApprox(equateApproxFrac, 0),
			}
			if d := cmp.Diff(tc.want, info, opts...); d != "" {
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

func unixTime(t *testing.T, year int, month time.Month, day int) int64 {
	t.Helper()
	var hh, mm, ss, ns int
	d := time.Date(year, month, day, hh, mm, ss, ns, time.UTC)
	return d.Unix()
}
