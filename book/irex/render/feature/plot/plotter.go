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
	tr Transformer
}

// NewPlotter creates a plotter to plot quotes in the virtual coordinates box
// from (0,0) to (1,1) and translate the coordinates using tr translator.
func NewPlotter(tr Transformer) *plotter {
	return &plotter{tr: tr}
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
		tr := roundTransformer{WithTransformer(Scale(1.0/2, 1.0/2), pl.tr)}
		x, y := tr.Transform(1, 1)
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: y,
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
	tr := roundTransformer{pl.initTransformer(quotes)}
	return pl.plot(quotes, tr)
}

func (pl *plotter) initTransformer(quotes []*pb.PlotFeature_Quote) Transformer {
	var ymin, ymax float64 = math.MaxFloat64, 0
	for _, q := range quotes {
		v := float64(q.GetCent().GetValue())
		if v > ymax {
			ymax = v
		}
		if v < ymin {
			ymin = v
		}
	}
	sx := 1 / float64(len(quotes)-1)
	sy := 1 / float64(ymax-ymin)
	return WithTransformer(Translate(0, -ymin), Scale(sx, sy), pl.tr)
}

func (pl *plotter) plot(quotes []*pb.PlotFeature_Quote, tr roundTransformer) (*Path, XQuote) {
	p := &Path{
		Commands: make([]PathCommand, len(quotes)),
	}
	qq := make(XQuote)
	for idx, q := range quotes {
		x, y := tr.Transform(float64(idx), float64(q.GetCent().GetValue()))
		if idx == 0 {
			p.Commands[idx] = PathMoveCommand{
				Point: Point{X: x, Y: y},
			}
		} else {
			p.Commands[idx] = PathLineCommand{
				Point: Point{X: x, Y: y},
			}
		}
		qq[x] = newQuote(q)
	}
	return p, qq
}

func round(x float64) int {
	return int(math.Floor(x))
}

type roundTransformer struct {
	t Transformer
}

// Transform applies transformation of (x,y) and then rounds the result.
func (rt roundTransformer) Transform(x, y float64) (int, int) {
	x, y = rt.t.Transform(x, y)
	return round(x), round(y)
}
