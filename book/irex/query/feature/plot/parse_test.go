// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query/feature/plot"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		q       string
		want    *pb.PlotIntent
		wantErr error
	}{
		{
			name:    "empty params",
			wantErr: plot.ErrPlotNoSymbol,
		},
		{
			name:    "invalid symbol",
			q:       "test-symbol",
			wantErr: plot.ErrPlotSymbol,
		},
		{
			name: "spx",
			q:    "spx",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_IDX_SPX)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := plot.Parse(tc.q)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%q unexpected error '%v'; want '%v'", tc.q, err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("%q mismatch (-want +got):\n%s", tc.q, d)
			}
		})
	}
}

func newPlotIntent(t *testing.T, s *pb.Symbol) *pb.PlotIntent {
	t.Helper()
	return pb.PlotIntent_builder{Symbol: s}.Build()
}

func newIndexSymbol(t *testing.T, idx pb.Symbol_Index) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{Index: idx.Enum()}.Build()
}
