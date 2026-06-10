// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"embed"
	"html/template"
	"strings"
	"unicode"

	"github.com/skhal/lab/book/irex/pb"
)

var (
	//go:embed static
	efs embed.FS

	tmplPlotFeature = template.Must(template.New("index.html").ParseFS(efs, "static/index.html"))
)

// renderer renders quores from the PlotFeature in SVG format. It adds x and y
// axes and uses the plotter to plot the quotes inside the svg view box.
type renderer struct {
	msg           *pb.PlotFeature
	width, height int
	axisOffset    int
}

const (
	defaultWidth      = 400
	defaultHeight     = 200
	defaultAxisOffset = 10
)

// NewRenderer creates a renderer for PlotFeature.
func NewRenderer(msg *pb.PlotFeature) *renderer {
	return &renderer{
		msg:        msg,
		width:      defaultWidth,
		height:     defaultHeight,
		axisOffset: defaultAxisOffset,
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
			Width:  fr.width,
			Height: fr.height,
		},
		Origin: &Point{
			X: fr.axisOffset,
			Y: fr.height - fr.axisOffset,
		},
		X: &Axis{
			Line: &Path{
				Move: Point{
					X: fr.axisOffset,
					Y: fr.height - fr.axisOffset,
				},
				Line: []Point{
					{
						X: fr.width - fr.axisOffset,
						Y: fr.height - fr.axisOffset,
					},
				},
			},
		},
		Y: &Axis{
			Line: &Path{
				Move: Point{
					X: fr.axisOffset,
					Y: fr.axisOffset,
				},
				Line: []Point{
					{
						X: fr.axisOffset,
						Y: fr.height - fr.axisOffset,
					},
				},
			},
		},
		Path: &Path{Line: line},
	}
}

func (fr *renderer) plot() []Point {
	xrange := fr.width - fr.axisOffset
	yrange := fr.height - fr.axisOffset
	pl := NewPlotter(xrange, yrange)
	return pl.Plot(fr.msg.GetQuotes())
}

// Axis describes a single axis.
type Axis struct {
	// Line draws the axis line.
	Line *Path
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
