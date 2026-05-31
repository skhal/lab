// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package main

import (
	"syscall/js"
)

type dimensions struct {
	width  int
	height int
}

type drawer struct {
	ctx  *canvasRenderingContext2D
	dim  dimensions
	sun  *astroObj
	done chan (struct{})
}

func newDrawer(ctx *canvasRenderingContext2D, dim dimensions) *drawer {
	const (
		sunRadius        = 432.2
		earthOrbitRadius = 93
		earthOrbitDays   = 365
		earthSpinHours   = 24
	)
	const planetRadiusScale = 50
	radius := func(r float64) float64 {
		// r is in thousand of kilometers
		const scale = 30
		return r / sunRadius * scale
	}
	orbitRadius := func(r float64) float64 {
		// r is in million miles
		const scale = sunRadius * 0.3
		return r / earthOrbitRadius * scale
	}
	orbitSeconds := func(d float64) float64 {
		// d is in days
		const scale = 40 // one orbit day in 1 second
		return d / earthOrbitDays * scale
	}
	spinSeconds := func(h float64) float64 {
		// d is in hours
		const scale = 5
		return h / earthSpinHours * scale
	}
	sun := newAstoObj(ctx, Config{
		Radius:      radius(432.2),
		SpinSeconds: spinSeconds(30 * 24),
		Fill:        &Color{R: 255, G: 255, B: 0},
	})
	mercury := newAstoObj(ctx, Config{
		Radius:       radius(1.5) * planetRadiusScale,
		OrbitRadius:  orbitRadius(35),
		OrbitSeconds: orbitSeconds(88),
		SpinSeconds:  spinSeconds(59),
		Fill:         &Color{R: 128, G: 128, B: 128},
	})
	sun.AddSatellite(mercury)
	venus := newAstoObj(ctx, Config{
		Radius:       radius(3.7) * planetRadiusScale,
		OrbitRadius:  orbitRadius(67),
		OrbitSeconds: orbitSeconds(224),
		SpinSeconds:  spinSeconds(243),
		Fill:         &Color{R: 192, G: 192, B: 192},
	})
	sun.AddSatellite(venus)
	earth := newAstoObj(ctx, Config{
		Radius:       radius(3.9) * planetRadiusScale,
		OrbitRadius:  orbitRadius(93),
		OrbitSeconds: orbitSeconds(365),
		SpinSeconds:  spinSeconds(24),
		Fill:         &Color{R: 42, G: 128, B: 0},
	})
	sun.AddSatellite(earth)
	mars := newAstoObj(ctx, Config{
		Radius:       radius(2.1) * planetRadiusScale,
		OrbitRadius:  orbitRadius(140),
		OrbitSeconds: orbitSeconds(687),
		SpinSeconds:  spinSeconds(24.4),
		Fill:         &Color{R: 255, G: 69, B: 0},
	})
	sun.AddSatellite(mars)
	return &drawer{ctx: ctx, dim: dim, sun: sun}
}

// Start starts the animation loop.
func (d *drawer) Start() {
	d.done = make(chan struct{})
	go func() {
		var f js.Func
		defer f.Release()
		var afreq int
		f = js.FuncOf(func(js.Value, []js.Value) any {
			select {
			case <-d.done:
				js.Global().Call("cancelAnimationFrame", afreq)
				return nil
			default:
			}
			d.render()
			afreq = js.Global().Call("requestAnimationFrame", f).Int()
			return nil
		})
		afreq = js.Global().Call("requestAnimationFrame", f).Int()
		<-d.done
	}()
}

// Stop stops the animation loop
func (d *drawer) Stop() {
	close(d.done)
}

func (d *drawer) render() {
	contextLock(d.ctx, d.reset)
	contextLock(d.ctx, d.renderSun)
}

func (d *drawer) reset() {
	d.ctx.FillStyle("black")
	d.ctx.FillRect(0, 0, d.dim.width, d.dim.height)
}

func (d *drawer) renderSun() {
	d.ctx.Translate(d.dim.width/2, d.dim.width/2)
	d.sun.Draw()
}

func contextLock(ctx *canvasRenderingContext2D, f func()) {
	ctx.Save()
	defer ctx.Restore()
	f()
}
