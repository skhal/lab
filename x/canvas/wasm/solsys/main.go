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

func main() {
	var app *application
	js.Global().Set("Run", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if app != nil {
			err := fmt.Errorf("application is already initialized")
			console().Error(err.Error())
			return nil
		}
		if len(args) < 1 {
			err := fmt.Errorf("Run needs canvas ID")
			console().Error(err.Error())
			return nil
		}
		idCanvas := args[0].String()
		var err error
		app, err = newApplication(idCanvas)
		if err != nil {
			console().Error(err.Error())
			return nil
		}
		return nil
	}))
	js.Global().Set("StartAnimation", js.FuncOf(func(js.Value, []js.Value) any {
		if err := app.StartAnimation(); err != nil {
			console().Error(err.Error())
		}
		return nil
	}))
	js.Global().Set("StopAnimation", js.FuncOf(func(js.Value, []js.Value) any {
		if err := app.StopAnimation(); err != nil {
			console().Error(err.Error())
		}
		return nil
	}))
	select {}
}

type application struct {
	d *drawer
}

func newApplication(cid string) (*application, error) {
	canvas, err := newCanvas(cid)
	if err != nil {
		return nil, err
	}
	ctx, err := newCanvasRenderingContext2D(canvas)
	if err != nil {
		return nil, err
	}
	dim := dimensions{width: canvas.width(), height: canvas.height()}
	return &application{
		d: newDrawer(ctx, dim),
	}, nil
}

// StartAnimation starts the animation of the canvas.
func (app *application) StartAnimation() error {
	app.d.Start()
	return nil
}

// StopAnimation stops the animation of the canvas.
func (app *application) StopAnimation() error {
	app.d.Stop()
	return nil
}
