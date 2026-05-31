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

type document struct {
	v js.Value
}

func newDocument() *document {
	return &document{
		v: js.Global().Get("document"),
	}
}

// GetElementByID retrieves a JavaScript object by identifier.
func (d *document) GetElementByID(id string) (js.Value, error) {
	v := d.v.Call("getElementById", id)
	switch {
	case v.IsNull(), v.IsUndefined():
		return js.Null(), fmt.Errorf("element #%s does not exist", id)
	}
	return v, nil
}
