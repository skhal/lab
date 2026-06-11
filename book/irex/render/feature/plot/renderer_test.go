// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/render/feature/plot"
	labtesting "github.com/skhal/lab/go/tests"
)

var update = flag.Bool("update", false, "update golden files")

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name       string
		symbol     *pb.Symbol
		quotes     []*pb.PlotFeature_Quote
		goldenHTML labtesting.GoldenFile
		goldenJS   labtesting.GoldenFile
		wantErr    error
	}{
		{
			name:   "one quote",
			symbol: newIndexSymbol(t, pb.Symbol_IDX_SPX),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
			},
			goldenHTML: labtesting.GoldenFile("testdata/one_quote.html"),
			goldenJS:   labtesting.GoldenFile("testdata/one_quote.js"),
		},
		{
			name:   "two quotes",
			symbol: newIndexSymbol(t, pb.Symbol_IDX_SPX),
			quotes: []*pb.PlotFeature_Quote{
				newQuote(t, 1990, time.January, 31, 101),
				newQuote(t, 1990, time.February, 28, 201),
			},
			goldenHTML: labtesting.GoldenFile("testdata/two_quotes.html"),
			goldenJS:   labtesting.GoldenFile("testdata/two_quotes.js"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := pb.PlotFeature_builder{
				Symbol: tc.symbol,
				Quotes: tc.quotes,
			}.Build()
			fr := plot.NewRenderer(msg)

			html, js, err := fr.Render()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
			if *update {
				tc.goldenHTML.Write(t, string(html))
				tc.goldenJS.Write(t, string(js))
			}
			if d := tc.goldenHTML.Diff(t, string(html)); d != "" {
				t.Errorf("HTML mismatch (-want +got):\n%s", d)
			}
			if d := tc.goldenJS.Diff(t, string(js)); d != "" {
				t.Errorf("JS mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func newIndexSymbol(t *testing.T, idx pb.Symbol_Index) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{Index: idx.Enum()}.Build()
}
