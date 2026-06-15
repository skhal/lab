// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/market/service"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestService_Quote(t *testing.T) {
	tests := []struct {
		name    string
		quotes  []*pb.Quote
		req     *pb.QuoteRequest
		wantRes *pb.QuoteResponse
		wantErr error
	}{
		{
			name: "missing symbol",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
			},
			req:     pb.QuoteRequest_builder{}.Build(),
			wantErr: service.ErrInvalidSymbol,
		},
		{
			name: "spx",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
			},
			req: pb.QuoteRequest_builder{Symbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX)}.Build(),
			wantRes: pb.QuoteResponse_builder{
				Quotes: []*pb.QuoteResponse_Quote{
					newResponseQuote(t, 1990, time.January, 31, 101),
					newResponseQuote(t, 1990, time.February, 28, 201),
				},
			}.Build(),
		},
		{
			name: "spx since date",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			req: pb.QuoteRequest_builder{
				Symbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
				Since:  newDate(t, 1990, time.February, 28),
			}.Build(),
			wantRes: pb.QuoteResponse_builder{
				Quotes: []*pb.QuoteResponse_Quote{
					newResponseQuote(t, 1990, time.February, 28, 201),
					newResponseQuote(t, 1990, time.March, 31, 301),
				},
			}.Build(),
		},
		{
			name: "spx until date",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			req: pb.QuoteRequest_builder{
				Symbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
				Until:  newDate(t, 1990, time.March, 31),
			}.Build(),
			wantRes: pb.QuoteResponse_builder{
				Quotes: []*pb.QuoteResponse_Quote{
					newResponseQuote(t, 1990, time.January, 31, 101),
					newResponseQuote(t, 1990, time.February, 28, 201),
				},
			}.Build(),
		},
		{
			name: "spx since and until date",
			quotes: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102, 103),
				newQuote(t, 1990, time.February, 28, 201, 202, 203),
				newQuote(t, 1990, time.March, 31, 301, 302, 303),
			},
			req: pb.QuoteRequest_builder{
				Symbol: newIndexSymbol(t, pb.Symbol_Index_ID_SPX),
				Since:  newDate(t, 1990, time.February, 28),
				Until:  newDate(t, 1990, time.March, 31),
			}.Build(),
			wantRes: pb.QuoteResponse_builder{
				Quotes: []*pb.QuoteResponse_Quote{
					newResponseQuote(t, 1990, time.February, 28, 201),
				},
			}.Build(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := service.Service{Quotes: tc.quotes}

			res, err := s.Quote(t.Context(), tc.req)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.wantRes, res, protocmp.Transform()); d != "" {
				t.Errorf("response mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx, div, cpi int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: newDate(t, year, month, day),
		Spx:  pb.Cent_builder{Value: &spx}.Build(),
		Div:  pb.Cent_builder{Value: &div}.Build(),
		Cpi:  pb.Cent_builder{Value: &cpi}.Build(),
	}.Build()
}

func newDate(t *testing.T, year int32, month time.Month, day int32) *pb.Date {
	t.Helper()
	return pb.Date_builder{
		Year:  &year,
		Month: new(int32(month)),
		Day:   &day,
	}.Build()
}

func newResponseQuote(t *testing.T, year int32, month time.Month, day int32, cent int32) *pb.QuoteResponse_Quote {
	t.Helper()
	return pb.QuoteResponse_Quote_builder{
		Date: newDate(t, year, month, day),
		Cent: pb.Cent_builder{Value: &cent}.Build(),
	}.Build()
}

func newIndexSymbol(t *testing.T, idx pb.Symbol_Index_ID) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{
		Index: pb.Symbol_Index_builder{
			Id: &idx,
		}.Build(),
	}.Build()
}
