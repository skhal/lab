// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

// Transformer performs cartesian coordinates transformation.
type Transformer interface {
	// Transform transforms coordinates (x,y) in some way.
	Transform(x, y float64) (newx, newy float64)
}

type translateTransformer struct {
	x0, y0 float64
}

// Transform performs linear transformation of coordinates (x,y) to
// (x+x0,y+y0).
func (t *translateTransformer) Transform(x, y float64) (float64, float64) {
	return x + t.x0, y + t.y0
}

type scaleTransformer struct {
	rx, ry float64
}

// Transform scales coordinates (x,y) to (rx*x,ry*y).
func (t *scaleTransformer) Transform(x, y float64) (float64, float64) {
	return x * t.rx, y * t.ry
}

// Translate creates a linear transformer to transform coordinates (x,y) to
// (x+dx,y+dy).
func Translate(dx, dy float64) *translateTransformer {
	return &translateTransformer{x0: dx, y0: dy}
}

// Scaley creates a transformer to scale coordinates (x,y) to (rx*x,ry*y).
func Scale(rx, ry float64) *scaleTransformer {
	return &scaleTransformer{rx: rx, ry: ry}
}

// WithTransformer creates a chain of transformers starting with t.
func WithTransformer(t Transformer, opts ...Transformer) Transformer {
	if len(opts) == 0 {
		return t
	}
	queue := make([]Transformer, 1+len(opts))
	queue[0] = t
	for i, t := range opts {
		queue[i+1] = t
	}
	return transformerQueue(queue)
}

type transformerQueue []Transformer

// Transform runs a sequential chain of transformers. See [WithTransformer].
func (q transformerQueue) Transform(x, y float64) (float64, float64) {
	for _, t := range q {
		x, y = t.Transform(x, y)
	}
	return x, y
}
