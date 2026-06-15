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

var (
	// ErrInvalidSymbols means the requested symbol is invalid or unsupported.
	ErrInvalidSymbol = errors.New("invalid symbol")

	// ErrInvalidIndex means the requested index is invalid or unsupported.
	ErrInvalidIndex = errors.New("invalid index")

	// ErrInvalidMetric means the requested metric is invalid or unsupported.
	ErrInvalidMetric = errors.New("invalid metric")
)

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
		return svc.quoteIndex(sym.GetIndex(), &DateRange{req.GetSince(), req.GetUntil()})
	case pb.Symbol_Market_case:
		return svc.quoteMarket(sym.GetMarket(), &DateRange{req.GetSince(), req.GetUntil()})
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSymbol, sym)
	}
}

func (svc *Service) quoteIndex(idx *pb.Symbol_Index, dr *DateRange) (*pb.QuoteResponse, error) {
	switch idx.GetId() {
	case pb.Symbol_Index_ID_SPX:
		return svc.quoteIndexSPX(idx, dr)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidIndex, idx)
	}
}

func (svc *Service) quoteIndexSPX(idx *pb.Symbol_Index, dr *DateRange) (*pb.QuoteResponse, error) {
	cents := (*pb.Quote).GetSpx
	if idx.HasMetric() {
		switch m := idx.GetMetric(); m {
		case pb.Symbol_Index_MET_DIV:
			cents = (*pb.Quote).GetDiv
		case pb.Symbol_Index_MET_EARN:
			cents = (*pb.Quote).GetEarn
		default:
			return nil, fmt.Errorf("%w: %s", ErrInvalidMetric, m)
		}
	}
	return svc.collectValues(cents, dr), nil
}

func (svc *Service) quoteMarket(sym *pb.Symbol_Market, dr *DateRange) (*pb.QuoteResponse, error) {
	var cents func(*pb.Quote) *pb.Cent
	switch m := sym.GetMetric(); m {
	case pb.Symbol_Market_MET_CPI:
		cents = (*pb.Quote).GetCpi
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidMetric, m)
	}
	return svc.collectValues(cents, dr), nil
}

func (svc *Service) collectValues(f func(*pb.Quote) *pb.Cent, dr *DateRange) *pb.QuoteResponse {
	var quotes []*pb.QuoteResponse_Quote
	for q := range dr.Quotes(svc.Quotes) {
		quotes = append(quotes, pb.QuoteResponse_Quote_builder{
			Date: q.GetDate(),
			Cent: f(q),
		}.Build())
	}
	res := pb.QuoteResponse_builder{
		Quotes: quotes,
	}.Build()
	return res
}
