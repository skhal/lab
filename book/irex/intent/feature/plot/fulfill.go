// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package plot fulfills the plot feature.
package plot

import (
	"errors"
	"fmt"

	"github.com/skhal/lab/book/irex/market"
	"github.com/skhal/lab/book/irex/pb"
)

// ErrPlotSymbol means the financial identifier is unsupported.
var ErrPlotSymbol = errors.New("unsupported symbol")

// Fulfill returns a list of quotes for the requested symbol. It returns an
// error if the requested symbol is not supported.
func Fulfill(msg *pb.PlotIntent) (*pb.PlotFeature, error) {
	switch sym := msg.GetSymbol(); sym.WhichSymbolOneof() {
	case pb.Symbol_Index_case:
		return fulfillIndexPlotIntent(sym, msg)
	default:
		return nil, fmt.Errorf("fulfill plot intent: %w: %s", ErrPlotSymbol, sym)
	}
}

func fulfillIndexPlotIntent(sym *pb.Symbol, msg *pb.PlotIntent) (*pb.PlotFeature, error) {
	switch idx := sym.GetIndex(); idx.GetId() {
	case pb.Symbol_Index_ID_SPX:
		return fulfillSPXPlot(sym, msg)
	default:
		err := fmt.Errorf("fulfill plot intent: %w: index %s", ErrPlotSymbol, idx)
		return nil, err
	}
}

func fulfillSPXPlot(sym *pb.Symbol, msg *pb.PlotIntent) (*pb.PlotFeature, error) {
	req := pb.QuoteRequest_builder{
		Symbol: sym,
		Since:  msg.GetSince(),
		Until:  msg.GetUntil(),
	}.Build()
	res, err := market.Quote(req)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	quotes := make([]*pb.PlotFeature_Quote, len(res.GetQuotes()))
	for i, q := range res.GetQuotes() {
		quotes[i] = pb.PlotFeature_Quote_builder{
			Date: q.GetDate(),
			Cent: q.GetCent(),
		}.Build()
	}
	feature := pb.PlotFeature_builder{
		Symbol: sym,
		Quotes: quotes,
	}.Build()
	return feature, nil
}
