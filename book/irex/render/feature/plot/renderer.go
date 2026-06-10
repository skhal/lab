// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"html/template"
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
func (fr *renderer) Render() (template.HTML, error) {
	d := fr.generateTemplateData()
	var b strings.Builder
	if err := tmplPlotFeature.Execute(&b, d); err != nil {
		return "", err
	}
	s := strings.TrimLeftFunc(b.String(), unicode.IsSpace)
	return template.HTML(s), nil
}

func (fr *renderer) generateTemplateData() *TemplateData {
	line := fr.plot()
	return &TemplateData{
		Title: fr.msg.GetSymbol().String(),
		ViewBox: &ViewBox{
			Width:  fr.cfg.ViewBox.Width,
			Height: fr.cfg.ViewBox.Height,
		},
		Origin: &Point{
			X: fr.cfg.Axis.Width,
			Y: fr.cfg.ViewBox.Height - fr.cfg.Axis.Width,
		},
		X:    fr.plotXaxis(&fr.cfg.Axis),
		Y:    fr.plotYaxis(&fr.cfg.Axis),
		Path: line,
	}
}

func (fr *renderer) plotXaxis(cfg *AxisConfig) *Axis {
	const guideCount = 4
	line := func() *Path {
		y := fr.cfg.ViewBox.Height + cfg.Offset
		return &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: cfg.Width,
						Y: y - cfg.Width + cfg.TickWidth,
					},
				},
				PathLineCommand{
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
				PathLineCommand{
					Point: Point{
						X: fr.cfg.ViewBox.Width,
						Y: y - cfg.Width + cfg.TickWidth,
					},
				},
			},
		}
	}
	guides := func() *Path {
		p := &Path{
			Commands: make([]PathCommand, 2*guideCount),
		}
		var x int
		for i := range guideCount {
			switch i {
			case 0:
				x = cfg.Width
			case guideCount - 1:
				x = fr.cfg.ViewBox.Width
			default:
				x = cfg.Width + (fr.cfg.ViewBox.Width-cfg.Width)/(guideCount-1)*i
			}
			p.Commands[i*2] = PathMoveCommand{
				Point: Point{
					X: x,
					Y: 0,
				},
			}
			p.Commands[i*2+1] = PathLineCommand{
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

func (fr *renderer) plotYaxis(cfg *AxisConfig) *Axis {
	const guideCount = 5
	line := func() *Path {
		x := cfg.Width - cfg.Offset
		return &Path{
			Commands: []PathCommand{
				PathMoveCommand{
					Point: Point{
						X: x - cfg.TickWidth,
						Y: 0,
					},
				},
				PathLineCommand{
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
				PathLineCommand{
					Point: Point{
						X: x - cfg.TickWidth,
						Y: fr.cfg.ViewBox.Height - cfg.Width,
					},
				},
			},
		}
	}
	guides := func() *Path {
		p := &Path{
			Commands: make([]PathCommand, 2*guideCount),
		}
		var y int
		for i := range guideCount {
			switch i {
			case 0:
				y = 0
			case guideCount - 1:
				y = fr.cfg.ViewBox.Height - cfg.Width
			default:
				y = (fr.cfg.ViewBox.Height - cfg.Width) / (guideCount - 1) * i
			}
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
	return &Axis{
		Line:   line(),
		Guides: guides(),
	}
}

func (fr *renderer) plot() *Path {
	xrange := fr.cfg.ViewBox.Width - fr.cfg.Axis.Width
	yrange := fr.cfg.ViewBox.Height - fr.cfg.Axis.Width
	pl := NewPlotter(xrange, yrange)
	return pl.Plot(fr.msg.GetQuotes())
}

// Axis describes a single axis.
type Axis struct {
	// Line draws the axis line.
	Line *Path

	// Guides is a set of the axis guide lines.
	Guides *Path
}

// TemplateData is the input data to HTML template.
type TemplateData struct {
	// ViewBox defines visiple part of the user space in SVG.
	ViewBox *ViewBox

	// Origin defines the location of the graph origin on the SVG canvas.
	Origin *Point

	// X defines x-axis.
	X *Axis

	// Y defines y-axis.
	Y *Axis

	// Path is the plotted line of the quotes.
	Path *Path

	// Title is the name of the plot.
	Title string
}
