// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport_test

import (
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestImport(t *testing.T) {
	tests := []struct {
		name      string
		csv       string
		skipLines int
		scanLines int
		want      []*pb.Quote
		wantErr   error
	}{
		{
			name: "empty",
		},
		{
			name: "data",
			csv: `
1990.01,1.01,1.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
		{
			name: "sorts data",
			csv: `
1990.02,2.01,2.02
1990.01,1.01,1.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
				newQuote(t, 1990, time.February, 28, 201, 202),
			},
		},
		{
			name: "same date",
			csv: `
1990.01,1.01,1.02
1990.01,2.01,2.02
`,
			wantErr: csvimport.ErrImport,
		},
		{
			name: "date gap same year",
			csv: `
1990.01,1.01,1.02
1990.03,3.01,3.02
`,
			wantErr: csvimport.ErrImport,
		},
		{
			name: "date gap next year",
			csv: `
1990.01,1.01,1.02
1991.01,2.01,2.02
`,
			wantErr: csvimport.ErrImport,
		},
		{
			name: "no gap dec jan",
			csv: `
1990.12,1.01,1.02
1991.01,2.01,2.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.December, 31, 101, 102),
				newQuote(t, 1991, time.January, 31, 201, 202),
			},
		},
		{
			name: "data error",
			csv: `
1990.01,abc,1.02
`,
			wantErr: csvimport.ErrImport,
		},
		{
			name: "skip lines no data",
			csv: `
1-date,spx,dividend
`,
			skipLines: 1,
		},
		{
			name: "skip lines",
			csv: `
1-date,spx,dividend
1990.01,1.01,1.02
`,
			skipLines: 1,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
		{
			name: "skip lines data error",
			csv: `
1-date,spx,dividend
1990.01,abc,1.02
`,
			skipLines: 1,
			wantErr:   csvimport.ErrImport,
		},
		{
			name: "scan lines",
			csv: `
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			scanLines: 1,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
		{
			name: "skip lines scan lines same",
			csv: `
1-date,spx,dividend
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			skipLines: 1,
			scanLines: 1, // includes skip lines
		},
		{
			name: "skip lines scan lines",
			csv: `
1-date,spx,dividend
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			skipLines: 1,
			scanLines: 2,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := csv.NewReader(strings.NewReader(tc.csv))
			opts := []csvimport.Opt{
				csvimport.WithSkipLines(tc.skipLines),
				csvimport.WithScanLines(tc.scanLines),
			}

			got, err := csvimport.Import(r, opts...)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Logf("csv:\n%s", tc.csv)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx, div int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: newDate(t, year, month, day),
		Spx:  pb.Cent_builder{Value: &spx}.Build(),
		Div:  pb.Cent_builder{Value: &div}.Build(),
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
