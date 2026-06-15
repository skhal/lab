// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"time"

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

// Quote is a data point on the graph. It has a date and the value.
type Quote struct {
	// UnixTime is the number of seconds since Jan 1, 1970 UTC.
	// See [time.Time.Unix].
	UnixTime int64

	// Cents is the quote value on the date.
	Cents int32
}

func newQuote(pbq *pb.PlotFeature_Quote) *Quote {
	q := &Quote{
		Cents: pbq.GetCent().GetValue(),
	}
	pbd := pbq.GetDate()
	var hh, mm, ss, ns int
	d := time.Date(int(pbd.GetYear()), time.Month(pbd.GetMonth()), int(pbd.GetDay()), hh, mm, ss, ns, time.UTC)
	q.UnixTime = d.Unix()
	return q
}

// XQuote is the quote for x coordinate of the line.
type XQuote map[int]*Quote

// Plot plots the quotes and returns a list of points representing the graph.
func (pl *plotter) Plot(quotes []*pb.PlotFeature_Quote) (*Path, XQuote) {
	switch len(quotes) {
	case 0:
		return nil, nil
	case 1:
		// place a single quote in the middle of the plot
		x := int(pl.xrange / 2)
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: int(pl.yrange / 2),
					},
				},
			},
		}
		q := quotes[0]
		qq := XQuote{
			x: newQuote(q),
		}
		return p, qq
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

func (pl *plotter) plot(quotes []*pb.PlotFeature_Quote) (*Path, XQuote) {
	p := &Path{
		Commands: make([]PathCommand, len(quotes)),
	}
	qq := make(XQuote)
	tr := WithTransformer(Translate(0, -pl.ymin), Scale(pl.xrange/float64(len(quotes)-1), pl.yrange/float64(pl.ymax-pl.ymin)))
	for idx, q := range quotes {
		x, y := tr.Transform(float64(idx), float64(q.GetCent().GetValue()))
		if idx == 0 {
			p.Commands[idx] = PathMoveCommand{
				Point: Point{X: round(x), Y: round(y)},
			}
		} else {
			p.Commands[idx] = PathLineCommand{
				Point: Point{X: round(x), Y: round(y)},
			}
		}
		qq[round(x)] = newQuote(q)
	}
	return p, qq
}

func round(x float64) int {
	return int(math.Floor(x))
}
