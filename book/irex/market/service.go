// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package market

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
	switch sym := req.GetSymbol(); sym {
	case pb.QuoteRequest_SYMBOL_SPX:
		return svc.quoteSPX(ctx, req)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSymbol, sym)
	}
}

func (svc *Service) quoteSPX(context.Context, *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	quotes := make([]*pb.QuoteResponse_Quote, len(svc.Quotes))
	for i, q := range svc.Quotes {
		quotes[i] = pb.QuoteResponse_Quote_builder{
			Date: q.GetDate(),
			Cent: q.GetSpx(),
		}.Build()
	}
	res := pb.QuoteResponse_builder{
		Quotes: quotes,
	}.Build()
	return res, nil
}
