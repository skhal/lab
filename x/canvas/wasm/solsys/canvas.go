// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

type canvas struct {
	v js.Value
}

func newCanvas(id string) (*canvas, error) {
	c, err := newDocument().GetElementByID(id)
	if err != nil {
		return nil, err
	}
	return &canvas{v: c}, nil
}

func (c *canvas) width() int {
	return c.v.Get("width").Int()
}

func (c *canvas) height() int {
	return c.v.Get("height").Int()
}

// GetContext retrieves canvas context: 2D rendering, WebGL, etc.
// It returns an error if the requested context is not available.
func (c *canvas) GetContext(s string) (js.Value, error) {
	v := c.v.Call("getContext", s)
	switch {
	case v.IsNull(), v.IsUndefined():
		return js.Null(), fmt.Errorf("can't get context %s", s)
	}
	return v, nil
}

// canvasRenderingContext2D is an API into JavaScript CanvasRenderingContext2D.
// https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
type canvasRenderingContext2D struct {
	v js.Value
}

func newCanvasRenderingContext2D(c *canvas) (*canvasRenderingContext2D, error) {
	v, err := c.GetContext("2d")
	if err != nil {
		return nil, err
	}
	return &canvasRenderingContext2D{v: v}, nil
}

// Arc draws an arc of a circle centered at (x,y) with radius r. The arc spans
// the angle (start,end) counter clock wise if ccw is false.
func (ctx *canvasRenderingContext2D) Arc(x, y int, r int, start, end float32, ccw bool) {
	ctx.v.Call("arc", x, y, r, start, end, ccw)
}

// BeginPath starts a path.
func (ctx *canvasRenderingContext2D) BeginPath() {
	ctx.v.Call("beginPath")
}

// Clip turns current path into a clipping region:
func (ctx *canvasRenderingContext2D) Clip() {
	ctx.v.Call("clip")
}

// canvasGradient implements HTML CanvasGradient.
// https://developer.mozilla.org/en-US/docs/Web/API/CanvasGradient
type canvasGradient struct {
	v js.Value
}

// AddColorStop adds a new color stop at offset and color value.
func (cg *canvasGradient) AddColorStop(x float32, c string) {
	cg.v.Call("addColorStop", x, c)
}

// CreateRadialGradient creates a radial gradient using the size and locations
// of two circles.
func (ctx *canvasRenderingContext2D) CreateRadialGradient(x1, y1, r1, x2, y2, r2 int) (*canvasGradient, error) {
	v := ctx.v.Call("createRadialGradient", x1, y1, r1, x2, y2, r2)
	switch {
	case v.IsNull(), v.IsUndefined():
		return nil, fmt.Errorf("failed to create radial gradient")
	}
	return &canvasGradient{v: v}, nil
}

// Fill fills the current path.
func (ctx *canvasRenderingContext2D) Fill() {
	ctx.v.Call("fill")
}

// FillStyle sets fill style, e.g. color.
func (ctx *canvasRenderingContext2D) FillStyle(s string) {
	ctx.v.Set("fillStyle", s)
}

// FillStyleGradient sets gradient fill style.
func (ctx *canvasRenderingContext2D) FillStyleGradient(cg *canvasGradient) {
	ctx.v.Set("fillStyle", cg.v)
}

// FillRect creates a filled rectangle at position (x,y) with size (w,h).
func (ctx *canvasRenderingContext2D) FillRect(x, y int, w, h int) {
	ctx.v.Call("fillRect", x, y, w, h)
}

// LineTo adds a line segment from current point to the point (x,y).
func (ctx *canvasRenderingContext2D) LineTo(x, y int) {
	ctx.v.Call("lineTo", x, y)
}

// LineWidth sets the width of the line to be used to strok paths.
func (ctx *canvasRenderingContext2D) LineWidth(w int) {
	ctx.v.Set("lineWidth", w)
}

// MoveTo sets the starting point of the new path to (x,y).
func (ctx *canvasRenderingContext2D) MoveTo(x, y int) {
	ctx.v.Call("moveTo", x, y)
}

// Restore loads last stored canvas state from the stack (see [Save]).
func (ctx *canvasRenderingContext2D) Restore() {
	ctx.v.Call("restore")
}

// Rotate turns the coordinates by angle phi in clock-wise direction.
func (ctx *canvasRenderingContext2D) Rotate(phi float32) {
	ctx.v.Call("rotate", phi)
}

// Save stores properties of the canvas state including transformation matrix,
// styles, etc. on the stack.
func (ctx *canvasRenderingContext2D) Save() {
	ctx.v.Call("save")
}

// SetLineDash sets current line dash pattern.
func (ctx *canvasRenderingContext2D) SetLineDash(a []int) {
	ctx.v.Call("setLineDash", a)
}

// Stroke draws current path using set stroke style.
func (ctx *canvasRenderingContext2D) Stroke() {
	ctx.v.Call("stroke")
}

// StrokeStyle sets color of the strokes around shapes.
func (ctx *canvasRenderingContext2D) StrokeStyle(s string) {
	ctx.v.Set("strokeStyle", s)
}

// Translate moves the origin of the coordinates to the point (x,y).
func (ctx *canvasRenderingContext2D) Translate(x, y int) {
	ctx.v.Call("translate", x, y)
}
