// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"html/template"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/skhal/lab/book/irex/pb"
)

// renderer renders quores from the PlotFeature in SVG format. It adds x and y
// axes and uses the plotter to plot the quotes inside the svg view box.
type renderer struct {
	msg *pb.PlotFeature
	cfg PlotConfig
}

// PlotConfig defines configuration for the plot in SVG araea.
type PlotConfig struct {
	// ViewBox is the configuration for SVG view box.
	ViewBox ViewBoxConfig

	// Axes holds configuration for x- and y- axes.
	Axes AxesConfig
}

// ViewBoxConfig defines the SVG's view box.
type ViewBoxConfig struct {
	// Width is the width of the view box.
	Width int

	// Height is the height of the view box.
	Height int
}

// AxesConfig configurations for x- and y- axes.
type AxesConfig struct {
	// X is the x-axis configuration.
	X AxisConfig

	// Y is the y-axis configuration.
	Y AxisConfig
}

// AxisConfig configures an axis.
type AxisConfig struct {
	// X is the x-coordinate of the axis
	X int

	// Y is the y-coordinate of the axis
	Y int

	// Width of the axis, excluding the offset.
	Width int

	// Height of the axis, excluding the offset.
	Height int

	// Padding is the axis offset from the plot to make axes stand out.
	Padding int

	// TickWidth is the tick width of the axis.
	TickWidth int
}

// NewRenderer creates a renderer for PlotFeature.
func NewRenderer(msg *pb.PlotFeature) *renderer {
	vb := ViewBoxConfig{
		Width:  400,
		Height: 200,
	}
	xaxis := AxisConfig{
		Height:    30,
		Padding:   3,
		TickWidth: 5,
	}
	yaxis := AxisConfig{
		Width:     50,
		Height:    vb.Height - xaxis.Height,
		Padding:   3,
		TickWidth: 5,
	}
	xaxis.X = yaxis.Width
	xaxis.Y = vb.Height - xaxis.Height
	xaxis.Width = vb.Width - yaxis.Width
	return &renderer{
		msg: msg,
		cfg: PlotConfig{
			ViewBox: vb,
			Axes: AxesConfig{
				X: xaxis,
				Y: yaxis,
			},
		},
	}
}

// Render plots the feature and returns an SVG with plot.
// It returns an error if rendering fails.
func (fr *renderer) Render() (template.HTML, template.JS, error) {
	d := fr.generateTemplateData()
	return fr.executeTemplates(d)
}

func (fr *renderer) executeTemplates(d *TemplateData) (template.HTML, template.JS, error) {
	var b strings.Builder
	if err := tmplPlotFeature.Execute(&b, d.html); err != nil {
		return "", "", err
	}
	html := strings.TrimLeftFunc(b.String(), unicode.IsSpace)
	b.Reset()
	if err := tmplsPlotFeatureJS.Execute(&b, d.js); err != nil {
		return "", "", err
	}
	mjs := b.String()
	return template.HTML(html), template.JS(mjs), nil
}

func (fr *renderer) generateTemplateData() *TemplateData {
	info := fr.plot()
	return &TemplateData{
		html: &HTMLTemplateData{
			Title: fr.msg.GetSymbol().String(),
			ViewBox: &ViewBox{
				Width:  fr.cfg.ViewBox.Width,
				Height: fr.cfg.ViewBox.Height,
			},
			X:    fr.plotXaxis(&fr.cfg.Axes.X, info.Quotes),
			Y:    fr.plotYaxis(&fr.cfg.Axes.Y, info.Ymin, info.Ymax),
			Path: info.Path,
			Graph: &Graph{
				X:      fr.cfg.Axes.Y.Width,
				Y:      0,
				Width:  fr.cfg.ViewBox.Width - fr.cfg.Axes.Y.Width,
				Height: fr.cfg.ViewBox.Height - fr.cfg.Axes.X.Height,
			},
		},
		js: &JSTemplateData{
			Quotes: info.Quotes,
		},
	}
}

func (fr *renderer) plotXaxis(cfg *AxisConfig, quotes XQuote) *Axis {
	guideCount := func() int {
		for n := range 5 {
			if n < 3 {
				continue
			}
			if len(quotes)%n == 0 {
				return n
			}
		}
		return 2
	}()
	guidePositions := func() []int {
		step := cfg.Width / (guideCount - 1)
		xx := make([]int, guideCount)
		position := func(i int) int {
			switch i {
			case 0:
				return cfg.X
			case guideCount - 1:
				return cfg.X + cfg.Width
			default:
				return cfg.X + i*step
			}
		}
		for i := range guideCount {
			xx[i] = position(i)
		}
		return xx
	}()
	line := func() *Path {
		y := cfg.Y + cfg.Padding
		// main line
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: cfg.X,
						Y: y,
					},
				},
				PathLineCommand{
					Point: Point{
						X: cfg.X + cfg.Width,
						Y: y,
					},
				},
			},
		}
		// guides
		for i, x := range guidePositions {
			tickWidth := cfg.TickWidth
			if i > 0 && i+1 < guideCount {
				tickWidth /= 2
			}
			p.Commands = append(p.Commands,
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: y,
					},
				},
				PathLineCommand{
					Point: Point{
						X: x,
						Y: y + tickWidth,
					},
				})
		}
		return p
	}
	guides := func() *Path {
		p := &Path{
			Commands: make([]PathCommand, 2*guideCount),
		}
		for i, x := range guidePositions {
			p.Commands[2*i] = PathMoveCommand{
				Point: Point{
					X: x,
					Y: 0,
				},
			}
			p.Commands[2*i+1] = PathLineCommand{
				Point: Point{
					X: x,
					Y: cfg.Y,
				},
			}
		}
		return p
	}
	labels := func() []Text {
		ll := make([]Text, guideCount)
		fontHeight := 13
		y := cfg.Y + cfg.Padding + cfg.TickWidth + fontHeight
		for i, x := range guidePositions {
			quote, ok := quotes[x]
			if !ok {
				fmt.Printf("label at %d: failed to get the quote\n", x)
				continue
			}
			t := time.Unix(quote.UnixTime, 0)
			switch i {
			case 0:
			case guideCount - 1:
				// shift left
				x -= 50
			default:
				// shift left
				x -= 25
			}
			ll[i] = Text{
				X:   x,
				Y:   y,
				Val: t.Format("Jan 2006"),
			}
		}
		return ll
	}
	return &Axis{
		Line:   line(),
		Guides: guides(),
		Labels: labels(),
	}
}

func (fr *renderer) plotYaxis(cfg *AxisConfig, ymin, ymax float64) *Axis {
	const guideCount = 4
	guidePositions := func() []int {
		step := cfg.Height / (guideCount - 1)
		yy := make([]int, guideCount)
		position := func(i int) int {
			switch i {
			case 0:
				return cfg.Y
			case guideCount - 1:
				return cfg.Y + cfg.Height
			default:
				return cfg.Y + i*step
			}
		}
		for i := range guideCount {
			yy[i] = position(i)
		}
		return yy
	}()
	line := func() *Path {
		x := cfg.Width - cfg.Padding
		// main line
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: cfg.Y,
					},
				},
				PathLineCommand{
					Point: Point{
						X: x,
						Y: cfg.Y + cfg.Height,
					},
				},
			},
		}
		// guides
		for i, y := range guidePositions {
			tickWidth := cfg.TickWidth
			if i > 0 && i+1 < guideCount {
				tickWidth /= 2
			}
			p.Commands = append(p.Commands,
				PathMoveCommand{
					Point: Point{
						X: x - tickWidth,
						Y: y,
					},
				},
				PathLineCommand{
					Point: Point{
						X: x,
						Y: y,
					},
				})
		}
		return p
	}
	guides := func() *Path {
		p := &Path{
			Commands: make([]PathCommand, 2*guideCount),
		}
		for i, y := range guidePositions {
			p.Commands[2*i] = PathMoveCommand{
				Point: Point{
					X: cfg.Width,
					Y: y,
				},
			}
			p.Commands[2*i+1] = PathLineCommand{
				Point: Point{
					X: fr.cfg.ViewBox.Width,
					Y: y,
				},
			}
		}
		return p
	}
	labels := func() []Text {
		cents := func(n float64) string {
			n -= math.Mod(n, 100)
			n /= 100
			return fmt.Sprintf("%.0f", n)
		}
		ll := make([]Text, guideCount)
		const x = 15
		dy := (ymax - ymin) / float64(guideCount-1)
		var s string
		for i, y := range guidePositions {
			switch {
			case i == 0:
				y += 10 // shift down
				s = cents(ymax)
			case i+1 == guideCount:
				s = cents(ymin)
			default:
				y += 4 // shift down
				s = cents(ymax - float64(i)*dy)
			}
			ll[i] = Text{
				X:   x,
				Y:   y,
				Val: s,
			}
		}
		return ll
	}
	return &Axis{
		Line:   line(),
		Guides: guides(),
		Labels: labels(),
	}
}

func (fr *renderer) plot() *PlotInfo {
	pl := NewPlotter(fr.newTransformer())
	return pl.Plot(fr.msg.GetQuotes())
}

// newTransformer creates a transformer to convert cartesian coordinates with
// the graph in the rectangle (0,0)-(1,1) to SVG rectangle with flipped y-axis.
//
// The SVG graph rectangle is the graph visible area without axis, padding,
// etc. It has the following dimensions:
//
//	x: [Axis.Width, ViewBox.Width]
//	y: [0, ViewBox.Height - Axis.Width]
//
// Keep in mind that y-axis runs top-down, i.e. it is flipped.
//
// newTransformer transforms the graph into the graph rectangle using two
// transformations:
//
//  1. flip the graph and scale it to the graph rectangle diemsions
//  2. translate the flipped graph to the left-bottom corner of the graph
//     rectangle.
func (fr *renderer) newTransformer() Transformer {
	// scale up the plot to the graph range and flip the y-axis
	scale := func() Transformer {
		xrange := float64(fr.cfg.ViewBox.Width - fr.cfg.Axes.Y.Width)
		yrange := float64(fr.cfg.ViewBox.Height - fr.cfg.Axes.X.Height)
		return Scale(xrange, -yrange)
	}
	// move the graph to the left-bottom corner of the graph range.
	translate := func() Transformer {
		dx := float64(fr.cfg.Axes.Y.Width)
		dy := float64(fr.cfg.ViewBox.Height - fr.cfg.Axes.X.Height)
		return Translate(dx, dy)
	}
	return WithTransformer(scale(), translate())
}
