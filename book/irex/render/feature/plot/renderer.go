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

	// Axis configures X and Y axes.
	Axis AxisConfig
}

// ViewBoxConfig defines the SVG's view box.
type ViewBoxConfig struct {
	// Width is the width of the view box.
	Width int

	// Height is the height of the view box.
	Height int
}

// AxisConfig configures an axis.
type AxisConfig struct {
	// Offset is the axis offset from the plot to make axes stand out.
	Offset int

	// Width is the width of the axis.
	Width int

	// TickWidth is the tick width of the axis.
	TickWidth int
}

// NewRenderer creates a renderer for PlotFeature.
func NewRenderer(msg *pb.PlotFeature) *renderer {
	return &renderer{
		msg: msg,
		cfg: PlotConfig{
			ViewBox: ViewBoxConfig{
				Width:  400,
				Height: 200,
			},
			Axis: AxisConfig{
				Offset:    3,
				Width:     10,
				TickWidth: 5,
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
			X:    fr.plotXaxis(&fr.cfg.Axis),
			Y:    fr.plotYaxis(&fr.cfg.Axis, info.Ymin, info.Ymax),
			Path: info.Path,
		},
		js: &JSTemplateData{
			Quotes: info.Quotes,
		},
	}
}

func (fr *renderer) plotXaxis(cfg *AxisConfig) *Axis {
	const guideCount = 3
	guidePositions := func() []int {
		xx := make([]int, guideCount)
		position := func(i int) int {
			switch i {
			case 0:
				return cfg.Width
			case guideCount - 1:
				return fr.cfg.ViewBox.Width
			default:
				return cfg.Width + (fr.cfg.ViewBox.Width-cfg.Width)/(guideCount-1)*i
			}
		}
		for i := range guideCount {
			xx[i] = position(i)
		}
		return xx
	}()
	line := func() *Path {
		y := fr.cfg.ViewBox.Height + cfg.Offset
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: cfg.Width,
						Y: y - cfg.Width,
					},
				},
				PathLineCommand{
					Point: Point{
						X: fr.cfg.ViewBox.Width,
						Y: y - cfg.Width,
					},
				},
			},
		}
		for i, x := range guidePositions {
			tickWidth := cfg.TickWidth
			if i > 0 && i+1 < len(guidePositions) {
				tickWidth /= 2
			}
			p.Commands = append(p.Commands,
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: y - cfg.Width,
					},
				},
				PathLineCommand{
					Point: Point{
						X: x,
						Y: y - cfg.Width + tickWidth,
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
					Y: fr.cfg.ViewBox.Height - cfg.Width,
				},
			}
		}
		return p
	}
	return &Axis{
		Line:   line(),
		Guides: guides(),
	}
}

func (fr *renderer) plotYaxis(cfg *AxisConfig, ymin, ymax float64) *Axis {
	const guideCount = 4
	guidePositions := func() []int {
		yy := make([]int, guideCount)
		position := func(i int) int {
			switch i {
			case 0:
				return 0
			case guideCount - 1:
				return fr.cfg.ViewBox.Height - cfg.Width
			default:
				return (fr.cfg.ViewBox.Height - cfg.Width) / (guideCount - 1) * i
			}
		}
		for i := range guideCount {
			yy[i] = position(i)
		}
		return yy
	}()
	line := func() *Path {
		x := cfg.Width - cfg.Offset
		p := &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: x,
						Y: 0,
					},
				},
				PathLineCommand{
					Point: Point{
						X: x,
						Y: fr.cfg.ViewBox.Height - cfg.Width,
					},
				},
			},
		}
		for i, y := range guidePositions {
			tickWidth := cfg.TickWidth
			if i > 0 && i+1 < len(guidePositions) {
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
		ll := make([]Text, len(guidePositions))
		const x = 15
		dy := (ymax - ymin) / float64(len(guidePositions)-1)
		var s string
		for i, y := range guidePositions {
			switch {
			case i == 0:
				y += 10 // shift down
				s = cents(ymax)
			case i+1 == len(guidePositions):
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
		xrange := float64(fr.cfg.ViewBox.Width - fr.cfg.Axis.Width)
		yrange := float64(fr.cfg.ViewBox.Height - fr.cfg.Axis.Width)
		return Scale(xrange, -yrange)
	}
	// move the graph to the left-bottom corner of the graph range.
	translate := func() Transformer {
		dx := float64(fr.cfg.Axis.Width)
		dy := float64(fr.cfg.ViewBox.Height - fr.cfg.Axis.Width)
		return Translate(dx, dy)
	}
	return WithTransformer(scale(), translate())
}
