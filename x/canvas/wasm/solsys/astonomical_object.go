// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package main

import (
	"fmt"
	"math"
)

// the value may change, e.g. MacOS reduces the rate from 60 fps to 30 fps when
// running in "battery save" mode
const animationFPS = 30

var defaultFill = Color{R: 200, G: 200, B: 200, T: 50}

// Config describes an astronomical object.
type Config struct {
	Radius       float64 // radius of the object
	OrbitRadius  float64 // radius of the object's orbit
	OrbitSeconds float64 // seconds it takes to complete an orbit
	SpinSeconds  float64 // seconds it takes to complete a spin
	Fill         *Color  // fill color for the object
}

type astroObjOrbit struct {
	radius int
	dphi   float32
}

type astroObjSpin struct {
	dphi float32
}

type astroObjConfig struct {
	radius int
	orbit  astroObjOrbit
	spin   astroObjSpin
	fill   *Color
}

// Color describes RGB web color with optional transparency.
type Color struct {
	// R is the red component of the color.
	R int
	// G is the green component of the color.
	G int
	// B is the blue component of the color.
	B int
	// T is color transparency, must be non-zero to take effect.
	T int
}

// String dumps colors as "rgb(R G B)" string. It includes transparency index
// if non-zero.
func (c *Color) String() string {
	if c.T == 0 {
		return fmt.Sprintf("rgb(%d %d %d)", c.R, c.G, c.B)
	}
	return fmt.Sprintf("rgb(%d %d %d / %d%%)", c.R, c.G, c.B, c.T)
}

type astroObj struct {
	ctx *canvasRenderingContext2D
	cfg astroObjConfig

	frame      int
	satellites []*astroObj
}

func newAstoObj(ctx *canvasRenderingContext2D, cfg Config) *astroObj {
	c := astroObjConfig{
		radius: int(math.Floor(cfg.Radius)),
		orbit: astroObjOrbit{
			radius: int(math.Floor(cfg.OrbitRadius)),
			dphi:   dphi(cfg.OrbitSeconds),
		},
		spin: astroObjSpin{
			dphi: dphi(cfg.SpinSeconds),
		},
		fill: cfg.Fill,
	}
	if c.fill == nil {
		c.fill = new(defaultFill)
	}
	return &astroObj{ctx: ctx, cfg: c}
}

func dphi(seconds float64) float32 {
	var frames = float32(seconds) * animationFPS
	return 2 * math.Pi / frames
}

// AddSatellite adds sat object as a dependency to the current object o. When
// o draws, it triggers drawing of satellites using o's local coordinates.
func (o *astroObj) AddSatellite(sat *astroObj) {
	o.satellites = append(o.satellites, sat)
}

// Draw renders astronomical object.
func (o *astroObj) Draw() {
	o.frame++
	if o.cfg.orbit.radius != 0 {
		contextLock(o.ctx, o.drawOrbit)
	}
	o.position()
	o.spin()
	contextLock(o.ctx, o.draw)
	for _, s := range o.satellites {
		contextLock(o.ctx, s.Draw)
	}
}

func (o *astroObj) position() {
	var angle = float32(o.frame) * o.cfg.orbit.dphi
	o.ctx.Rotate(angle)
	o.ctx.Translate(o.cfg.orbit.radius, 0)
}

func (o *astroObj) spin() {
	var angle = float32(o.frame) * o.cfg.spin.dphi
	o.ctx.Rotate(angle)
}

func (o *astroObj) draw() {
	scaleRadius := func(r int, scale float64) int {
		return int(math.Floor(float64(r) * scale))
	}
	g, err := o.ctx.CreateRadialGradient(0, 0, scaleRadius(o.cfg.radius, 0.3), 0, 0, scaleRadius(o.cfg.radius, 0.9))
	c := *o.cfg.fill
	if c.T == 0 {
		c.T = 100
	}
	g.AddColorStop(0.1, c.String())
	c.T = int(math.Floor(float64(c.T) * 0.8))
	g.AddColorStop(0.9, c.String())
	g.AddColorStop(1.0, "transparent")
	if err != nil {
		panic(err)
	}
	o.ctx.FillStyleGradient(g)
	o.ctx.BeginPath()
	o.ctx.MoveTo(o.cfg.radius, 0)
	o.ctx.Arc(0, 0, o.cfg.radius, 0, 2*math.Pi, false)
	o.ctx.Fill()

	o.ctx.StrokeStyle("black")
	o.ctx.LineWidth(2)
	o.ctx.BeginPath()
	o.ctx.MoveTo(0, 0)
	o.ctx.LineTo(o.cfg.radius, 0)
	o.ctx.Stroke()
}

func (o *astroObj) drawOrbit() {
	o.ctx.StrokeStyle("rgb(146 146 146/50%)")
	o.ctx.LineWidth(2)
	o.ctx.BeginPath()
	o.ctx.MoveTo(o.cfg.orbit.radius, 0)
	o.ctx.Arc(0, 0, o.cfg.orbit.radius, 0, 2*math.Pi, true)
	o.ctx.Stroke()
}
