// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/skhal/lab/book/irex/pb"
)

// ErrInvalidSymbols means the requested symbol is invalid or unsupported.
var ErrInvalidSymbol = errors.New("invalid symbol")

// Service implements MarketService
type Service struct {
	pb.UnimplementedMarketServiceServer

	// Quotes is date ordered market data.
	Quotes []*pb.Quote
}

// Quote implements MarketService.Quote RPC endpoint. It retrieves quotes for
// a requested symbol or returns an error if the symbol is invalid.
func (svc *Service) Quote(ctx context.Context, req *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	switch sym := req.GetSymbol(); sym.WhichSymbolOneof() {
	case pb.Symbol_Index_case:
		return svc.quoteIndex(sym, &DateRange{req.GetSince(), req.GetUntil()})
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSymbol, sym)
	}
}

func (svc *Service) quoteIndex(sym *pb.Symbol, dr *DateRange) (*pb.QuoteResponse, error) {
	switch idx := sym.GetIndex(); idx {
	case pb.Symbol_IDX_SPX:
		return svc.quoteIndexSPX(dr)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSymbol, sym)
	}
}

func (svc *Service) quoteIndexSPX(dr *DateRange) (*pb.QuoteResponse, error) {
	var quotes []*pb.QuoteResponse_Quote
	for q := range dr.Quotes(svc.Quotes) {
		quotes = append(quotes, pb.QuoteResponse_Quote_builder{
			Date: q.GetDate(),
			Cent: q.GetSpx(),
		}.Build())
	}
	res := pb.QuoteResponse_builder{
		Quotes: quotes,
	}.Build()
	return res, nil
}
