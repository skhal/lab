// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skhal/lab/book/irex/pb"
)

var (
	// ErrPlotNoSymbol means the plot command does not have a symbol parameter.
	ErrPlotNoSymbol = errors.New("missing symbol")

	// ErrPlotSymbol means the plot command uses unsupported symbol.
	ErrPlotSymbol = errors.New("invalid symbol")
)

var indexByName = map[string]pb.Symbol_Index{
	"spx": pb.Symbol_IDX_SPX,
}

// Parse parses plot-command parameters. It returns the intent upon successful
// parse or error if the parse fails.
func Parse(symbol string) (*pb.PlotIntent, error) {
	if symbol == "" {
		return nil, ErrPlotNoSymbol
	}
	name := strings.ToLower(symbol)
	idx, ok := indexByName[name]
	if !ok {
		return nil, fmt.Errorf("plot: %w: %s", ErrPlotSymbol, symbol)
	}
	sym := pb.Symbol_builder{Index: &idx}.Build()
	msg := pb.PlotIntent_builder{Symbol: sym}.Build()
	return msg, nil
}
