// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package query_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query"
	"github.com/skhal/lab/book/irex/query/feature/plot"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestUnderstand(t *testing.T) {
	tests := []struct {
		name    string
		q       string
		want    *pb.Intent
		wantErr error
	}{
		{
			name:    "empty query",
			wantErr: query.ErrNoCommand,
		},
		{
			name:    "unsupported command",
			q:       "unsupportedcommand abc123",
			wantErr: query.ErrInvalidCommand,
		},
		{
			name:    "plot command no params",
			q:       "plot",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name:    "plot invalid symbol",
			q:       "plot test-symbol",
			wantErr: plot.ErrNoSymbol,
		},
		{
			name: "plot spx",
			q:    "plot spx",
			want: newPlotIntent(t, newIndexSymbol(t, pb.Symbol_IDX_SPX)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := query.Understand(tc.q)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%q unexpected error '%v'; want '%v'", tc.q, err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("%q mismatch (-want +got):\n%s", tc.q, d)
			}
		})
	}
}

func newPlotIntent(t *testing.T, s *pb.Symbol) *pb.Intent {
	t.Helper()
	msg := pb.PlotIntent_builder{Symbol: s}.Build()
	intent := pb.Intent_builder{}.Build()
	proto.SetExtension(intent, pb.E_PlotIntent_PlotIntent, msg)
	return intent
}

func newIndexSymbol(t *testing.T, idx pb.Symbol_Index) *pb.Symbol {
	t.Helper()
	return pb.Symbol_builder{Index: idx.Enum()}.Build()
}
