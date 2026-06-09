// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"

	"github.com/skhal/lab/book/irex/pb"
)

type plotter struct {
	xrange, yrange float64
	ymin, ymax     float64
}

// NewPlotter creates a plotter to plot the quotes in the box from (0,0) to
// (xrange, yrange).
func NewPlotter(xrange, yrange int) *plotter {
	return &plotter{
		xrange: float64(xrange),
		yrange: float64(yrange),
	}
}

// Plot plots the quotes and returns a list of points representing the graph.
func (pl *plotter) Plot(quotes []*pb.PlotFeature_Quote) []Point {
	switch len(quotes) {
	case 0:
		return nil
	case 1:
		// place a single quote in the middle of the plot
		return []Point{{int(pl.xrange / 2), int(pl.yrange / 2)}}
	}
	pl.initAxis(quotes)
	return pl.plot(quotes)
}

func (pl *plotter) initAxis(quotes []*pb.PlotFeature_Quote) {
	pl.ymin, pl.ymax = math.MaxFloat64, 0
	for _, q := range quotes {
		v := float64(q.GetCent().GetValue())
		if v > pl.ymax {
			pl.ymax = v
		}
		if v < pl.ymin {
			pl.ymin = v
		}
	}
}

func (pl *plotter) plot(quotes []*pb.PlotFeature_Quote) []Point {
	xtr := NewTransformer(0, pl.xrange/float64(len(quotes)-1))
	ytr := NewTransformer(pl.ymin, pl.yrange/float64(pl.ymax-pl.ymin))
	path := make([]Point, len(quotes))
	for idx, q := range quotes {
		x := xtr.Transform(float64(idx))
		y := ytr.Transform(float64(q.GetCent().GetValue()))
		path[idx] = Point{round(x), round(y)}
	}
	return path
}

func round(x float64) int {
	return int(math.Floor(x))
}
